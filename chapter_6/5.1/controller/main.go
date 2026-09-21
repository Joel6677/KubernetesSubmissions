package main

import (
	"context"
	"fmt"
	"io"
	"net/http"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var gvr = schema.GroupVersionResource{
	Group:    "stable.dwk",
	Version:  "v1",
	Resource: "dummysites",
}

func download(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("website returned HTTP %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func createConfigMap(
	clientset *kubernetes.Clientset,
	namespace string,
	name string,
	html string,
) error {
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: name + "-html",
		},
		Data: map[string]string{
			"index.html": html,
		},
	}

	_, err := clientset.CoreV1().
		ConfigMaps(namespace).
		Create(context.Background(), configMap, metav1.CreateOptions{})

	if apierrors.IsAlreadyExists(err) {
		return nil
	}

	return err
}

func int32Ptr(i int32) *int32 {
	return &i
}

func createDeployment(
	clientset *kubernetes.Clientset,
	namespace string,
	name string,
) error {
	labels := map[string]string{
		"app": name,
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},

		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),

			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},

			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},

				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: "nginx:alpine",

							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 80,
								},
							},

							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "html",
									MountPath: "/usr/share/nginx/html",
								},
							},
						},
					},

					Volumes: []corev1.Volume{
						{
							Name: "html",

							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: name + "-html",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	_, err := clientset.AppsV1().
		Deployments(namespace).
		Create(context.Background(), deployment, metav1.CreateOptions{})

	if apierrors.IsAlreadyExists(err) {
		return nil
	}

	return err
}

func createService(
	clientset *kubernetes.Clientset,
	namespace string,
	name string,
) error {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},

		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": name,
			},

			Ports: []corev1.ServicePort{
				{
					Port: 80,
				},
			},
		},
	}

	_, err := clientset.CoreV1().
		Services(namespace).
		Create(context.Background(), service, metav1.CreateOptions{})

	if apierrors.IsAlreadyExists(err) {
		return nil
	}

	return err
}

func main() {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	watcher, err := dynamicClient.Resource(gvr).
		Namespace("").
		Watch(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	defer watcher.Stop()

	fmt.Println("watching for dummysites")

	for event := range watcher.ResultChan() {

		if event.Type != watch.Added {
			continue
		}

		site, ok := event.Object.(*unstructured.Unstructured)
		if !ok {
			fmt.Println("could not convert object to dummysite")
			continue
		}

		url, found, err := unstructured.NestedString(
			site.Object,
			"spec",
			"website_url",
		)

		if err != nil || !found || url == "" {
			fmt.Println("invalid website_url")
			continue
		}

		namespace := site.GetNamespace()
		name := site.GetName()

		fmt.Println("new dummysite:", name, url)

		html, err := download(url)
		if err != nil {
			fmt.Println("download failed:", err)
			continue
		}

		err = createConfigMap(clientset, namespace, name, html)
		if err != nil {
			fmt.Println("configmap failed:", err)
			continue
		}

		err = createDeployment(clientset, namespace, name)
		if err != nil {
			fmt.Println("deployment failed:", err)
			continue
		}

		err = createService(clientset, namespace, name)
		if err != nil {
			fmt.Println("service failed:", err)
			continue
		}

		fmt.Println("dummysite created successfully:", name)
	}
}
