/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true

// Gateway represents an instantiation of a service-traffic handling infrastructure.
type Gateway struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GatewaySpec   `json:"spec,omitempty"`
	Status GatewayStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// GatewayList contains a list of Gateway
type GatewayList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Gateway `json:"items"`
}

// GatewaySpec defines the desired state of Gateway.
//
// The Spec is split into two major pieces: listeners describing
// client-facing properties and routes that describe application-level
// routing.
//
// Not all possible combinations of options specified in the Spec are
// valid. Some invalid configurations can be caught synchronously via a
// webhook, but there are many cases that will require asynchronous
// signaling via the GatewayStatus block.
type GatewaySpec struct {
	// Class used for this Gateway. This is the name of a GatewayClass resource.
	Class string `json:"class"`
	// Listeners associated with this Gateway. Listeners define what addresses,
	// ports, protocols are bound on this Gateway.
	Listeners []Listener `json:"listeners"`

	// Servers binds virtual servers that are hosted by the gateway
	// to listeners.
	//
	// An implementation must validate that Dedicated listeners have no
	// more than 1 server binding.
	//
	// An implementation must validate that all the servers bound to the
	// same combined listener have a compatible discriminator.
	// TODO(jpeach): describe consequences of this validating failing.
	Servers []ServerBinding `json:"servers"`
}

// ServerBinding specifies which listener a virtual server should
// be attached to. A server can be bound to 1 or more listeners.
type ServerBinding struct {
	// ListenerName is the unique name of a listener.
	//
	// TODO(jpeach): We could let an empty name mean all listeners, but
	// there are trade-offs to that. For now this field must not be empty.
	//
	// +required
	ListenerName `json:"listenerName"`
	// Server is a reference to the virtual server to be attached.
	//
	// +required
	Server ServerObjectReference `json:"server"`
}

type ListenerProtocolType string

const ListenerProtocolTCP ListenerProtocolType = "TCP"
const ListenerProtocolUDP ListenerProtocolType = "UDP"

type Listenertype string

// DedicatedListenerType describes a listener with exactly on attached server.
// All traffic received by this listener is terminated by this server.
const DedicatedListenerType ListenerType = "Dedicated"

// CombinedListenerType describes a listener what may have multiple attached
// servers. The protocol recognized by this listener must have a discriminator
// that can be used to select the appropriate server for the incoming request.
// In the case of a streaming data over TLS, the ALPN protocol name can sepect
// between different StreamServer specification. In the case of HTTP traffic,
// the HTTP Host (i.e. H2 authority) provides the discriminator.
const CombinedListenerType ListenerType = "Combined"

// Listener defines a
type Listener struct {
	// Type specifies the server binding semantics of this listener.
	//
	// +required
	Type ListenerType `json:"type"`

	// Name is the listener's name and should be specified as an
	// RFC 1035 DNS_LABEL [1]:
	//
	// [1] https://tools.ietf.org/html/rfc1035
	//
	// Each listener of a Gateway must have a unique name. Name is used
	// for associating a listener in Gateway status.
	//
	// Support: Core
	//
	// +required
	Name string `json:"name"`
	// Address requested for this listener. This is optional and behavior
	// can depend on GatewayClass. If a value is set in the spec and
	// the request address is invalid, the GatewayClass MUST indicate
	// this in the associated entry in GatewayStatus.Listeners.
	//
	// Support:
	//
	// +optional
	Address *ListenerAddress `json:"address,omitempty"`
	// Port is a list of ports associated with the Address.
	//
	// Support:
	// +optional
	Port *int32 `json:"port,omitempty"`
	// Protocol is the network layer protocol this listener accepts.
	//
	// Support:
	// +optional
	Protocol ListenerProtocolType `json:"protocol,omitempty"`
	// Extension for this Listener.  The resource may be "configmap" (use
	// the empty string for the group) or an implementation-defined resource
	// (for example, resource "mylistener" in group "networking.acme.io").
	//
	// Support: custom.
	// +optional
	Extension *ListenerExtensionObjectReference `json:"extension,omitempty"`
}

const (
	// IPAddress is an address that is an IP address.
	//
	// Support: Extended.
	IPAddress = "IPAddress"
	// NamedAddress is an address selected by name. The interpretation of
	// the name is depenedent on the controller.
	//
	// Support: Implementation-specific.
	NamedAddress = "NamedAddress"
)

// ListenerAddress describes an address for the Listener.
type ListenerAddress struct {
	// Type of the Address. This is one of the *AddressType constants.
	//
	// Support: Extended
	Type string `json:"type"`
	// Value. Examples: "1.2.3.4", "128::1", "my-ip-address". Validity of the
	// values will depend on `Type` and support by the controller.
	Value string `json:"value"`
}

const (
	// TLS1_0 denotes the TLS v1.0.
	TLS1_0 = "TLS1_0"
	// TLS1_1 denotes the TLS v1.0.
	TLS1_1 = "TLS1_1"
	// TLS1_2 denotes the TLS v1.0.
	TLS1_2 = "TLS1_2"
	// TLS1_3 denotes the TLS v1.0.
	TLS1_3 = "TLS1_3"
)

// LocalObjectReference identifies an API object within a known namespace.
type LocalObjectReference struct {
	// Group is the group of the referent.  The empty string represents
	// the core API group.
	//
	// +kubebuilder:validation:Required
	// +required
	Group string `json:"group"`
	// Resource is the resource of the referent.
	//
	// +kubebuilder:validation:Required
	// +required
	Resource string `json:"resource"`
	// Name is the name of the referent.
	//
	// +kubebuilder:validation:Required
	// +required
	Name string `json:"name"`
}

// CertificateObjectReference identifies a certificate object within a known
// namespace.
//
// +k8s:deepcopy-gen=false
type CertificateObjectReference = LocalObjectReference

// ListenerExtensionObjectReference identifies a listener extension object
// within a known namespace.
//
// +k8s:deepcopy-gen=false
type ListenerExtensionObjectReference = LocalObjectReference

// ServerObjectReference identifies a virtual server object in the same
// namespace as the Gateway.
//
// +k8s:deepcopy-gen=false
type ServerObjectReference = LocalObjectReference

// GatewayStatus defines the observed state of Gateway.
type GatewayStatus struct {
	// Conditions describe the current conditions of the Gateway.
	Conditions []GatewayCondition `json:"conditions"`
	// Listeners provide status for each listener defined in the Spec. The name
	// in ListenerStatus refers to the corresponding Listener of the same name.
	Listeners []ListenerStatus `json:"listeners"`
}

// GatewayConditionType is a type of condition associated with a Gateway.
type GatewayConditionType string

const (
	// ConditionNoSuchGatewayClass indicates that the specified GatewayClass
	// does not exist.
	ConditionNoSuchGatewayClass GatewayConditionType = "NoSuchGatewayClass"
	// ConditionGatewayNotScheduled indicates that the Gateway has not been
	// scheduled.
	ConditionGatewayNotScheduled GatewayConditionType = "GatewayNotScheduled"
	// ConditionListenersNotReady indicates that at least one of the specified
	// listeners is not ready. If this condition has a status of True, a more
	// detailed ListenerCondition should be present in the corresponding
	// ListenerStatus.
	ConditionListenersNotReady GatewayConditionType = "ListenersNotReady"
	// ConditionInvalidListeners indicates that at least one of the specified
	// listeners is invalid. If this condition has a status of True, a more
	// detailed ListenerCondition should be present in the corresponding
	// ListenerStatus.
	ConditionInvalidListeners GatewayConditionType = "InvalidListeners"
	// ConditionRoutesNotReady indicates that at least one of the specified
	// routes is not ready.
	ConditionRoutesNotReady GatewayConditionType = "RoutesNotReady"
	// ConditionInvalidRoutes indicates that at least one of the specified
	// routes is invalid.
	ConditionInvalidRoutes GatewayConditionType = "InvalidRoutes"
)

// GatewayCondition is an error status for a given route.
type GatewayCondition struct {
	// Type indicates the type of condition.
	Type GatewayConditionType `json:"type"`
	// Status describes the current state of this condition. Can be "True",
	// "False", or "Unknown".
	Status core.ConditionStatus `json:"status"`
	// Message is a human-understandable message describing the condition.
	// +optional
	Message string `json:"message,omitempty"`
	// Reason indicates why the condition is in this state.
	// +optional
	Reason string `json:"reason,omitempty"`
	// LastTransitionTime indicates the last time this condition changed.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
}

// ListenerStatus is the status associated with each listener block.
type ListenerStatus struct {
	// Name is the name of the listener this status refers to.
	Name string `json:"name"`
	// Address bound on this listener.
	Address *ListenerAddress `json:"address"`
	// Conditions describe the current condition of this listener.
	Conditions []ListenerCondition `json:"conditions"`
}

// ListenerConditionType is a type of condition associated with the listener.
type ListenerConditionType string

const (
	// ConditionInvalidListener is a generic condition that is a catch all for
	// unsupported configurations that do not match a more specific condition.
	// Implementors should try to use a more specific condition instead of this
	// one to give users and automation more information.
	ConditionInvalidListener ListenerConditionType = "InvalidListener"
	// ConditionListenerNotReady indicates the listener is not ready.
	ConditionListenerNotReady ListenerConditionType = "ListenerNotReady"
	// ConditionInvalidAddress indicates the Address is invalid.
	ConditionInvalidAddress ListenerConditionType = "InvalidAddress"
)

// ListenerCondition is an error status for a given listener.
type ListenerCondition struct {
	// Type indicates the type of condition.
	Type ListenerConditionType `json:"type"`
	// Status describes the current state of this condition. Can be "True",
	// "False", or "Unknown".
	Status core.ConditionStatus `json:"status"`
	// Message is a human-understandable message describing the condition.
	// +optional
	Message string `json:"message,omitempty"`
	// Reason indicates why the condition is in this state.
	// +optional
	Reason string `json:"reason,omitempty"`
	// LastTransitionTime indicates the last time this condition changed.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
}

func init() {
	SchemeBuilder.Register(&Gateway{}, &GatewayList{})
}
