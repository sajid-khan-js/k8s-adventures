package main

// Pod model info
// @Description Pod information
type Pod struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

// Namespace model info
// @Description Namespace information
type Namespace struct {
	// https://github.com/go-playground/validator and https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-subdomain-names
	Name string `json:"name" binding:"required,hostname_rfc1123"`
	Pods []Pod  `json:"pods"`
}

/*

Test data:

var namespaces = []Namespace{
	{Name: "default", Pods: []Pod{
		Pod{Name: "nginx", Status: "Running"},
		Pod{Name: "httpbin", Status: "Pending"}}},
	{Name: "kube-system", Pods: []Pod{
		Pod{Name: "coredns-558bd4d5db-gmbdd", Status: "Running"},
		Pod{Name: "etcd-docker-desktop", Status: "Running"},
		Pod{Name: "kube-scheduler-docker-desktop", Status: "Running"}}},
	{Name: "app", Pods: []Pod{
		Pod{Name: "my-app", Status: "CrashLoopBackOff"}}},
}

*/
