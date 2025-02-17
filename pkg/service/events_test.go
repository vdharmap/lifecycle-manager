package service

import (
	"fmt"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func Test_NewEvent(t *testing.T) {
	t.Log("Test_NewEvent: should be able to get a new kubernetes event")
	referencedNode := "ip-10-10-10-10"
	msg := fmt.Sprintf(EventMessageInstanceDeregisterFailed, "i-123456789012", "my-load-balancer", "some bad error occured")
	event := newKubernetesEvent(EventReasonInstanceDeregisterFailed, msg, referencedNode)

	if event.Reason != string(EventReasonInstanceDeregisterFailed) {
		t.Fatalf("expected event.Reason to be: %v, got: %v", string(EventReasonInstanceDeregisterFailed), event.Reason)
	}

	if event.InvolvedObject.Name != referencedNode {
		t.Fatalf("expected event.InvolvedObject.Name to be: %v, got: %v", referencedNode, event.InvolvedObject.Name)
	}

	if event.Message != msg {
		t.Fatalf("expected event.Message to be: %v, got: %v", msg, event.Message)
	}

}

func Test_PublishEvent(t *testing.T) {
	t.Log("Test_PublishEvent: should be able to publish kubernetes events")
	kubeClient := fake.NewSimpleClientset()
	referencedNode := "ip-10-10-10-10"
	msg := fmt.Sprintf(EventMessageInstanceDeregisterFailed, "i-123456789012", "my-load-balancer", "some bad error occured")
	event := newKubernetesEvent(EventReasonInstanceDeregisterFailed, msg, referencedNode)

	publishKubernetesEvent(kubeClient, event)
	expectedEvents := 1

	events, err := kubeClient.CoreV1().Events(EventNamespace).List(metav1.ListOptions{})
	if err != nil {
		t.Fatalf("Test_PublishEvent: expected error not to have occured, %v", err)
	}

	if len(events.Items) != expectedEvents {
		t.Fatalf("Test_PublishEvent: expected %v events, found: %v", expectedEvents, len(events.Items))
	}
}
