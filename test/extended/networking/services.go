package networking

import (
	admissionapi "k8s.io/pod-security-admission/api"

	e2e "k8s.io/kubernetes/test/e2e/framework"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	exutil "github.com/openshift/origin/test/extended/util"
)

var _ = Describe("[sig-network] services", func() {
	Context("basic functionality", func() {
		f1 := e2e.NewDefaultFramework("net-services1")
		// TODO(sur): verify if privileged is really necessary in a follow-up
		f1.NamespacePodSecurityEnforceLevel = admissionapi.LevelPrivileged

		It("should allow connections to another pod on the same node via a service IP", func() {
			Expect(checkServiceConnectivity(f1, f1, SAME_NODE)).To(Succeed())
		})

		It("should allow connections to another pod on a different node via a service IP", func() {
			Expect(checkServiceConnectivity(f1, f1, DIFFERENT_NODE)).To(Succeed())
		})
	})

	InNonIsolatingContext(func() {
		f1 := e2e.NewDefaultFramework("net-services1")
		// TODO(sur): verify if privileged is really necessary in a follow-up
		f1.NamespacePodSecurityEnforceLevel = admissionapi.LevelPrivileged
		f2 := e2e.NewDefaultFramework("net-services2")
		// TODO(sur): verify if privileged is really necessary in a follow-up
		f2.NamespacePodSecurityEnforceLevel = admissionapi.LevelPrivileged

		It("should allow connections to pods in different namespaces on the same node via service IPs", func() {
			Expect(checkServiceConnectivity(f1, f2, SAME_NODE)).To(Succeed())
		})

		It("should allow connections to pods in different namespaces on different nodes via service IPs", func() {
			Expect(checkServiceConnectivity(f1, f2, DIFFERENT_NODE)).To(Succeed())
		})
	})

	oc := exutil.NewCLI("ns-global")

	InIsolatingContext(func() {
		f1 := e2e.NewDefaultFramework("net-services1")
		// TODO(sur): verify if privileged is really necessary in a follow-up
		f1.NamespacePodSecurityEnforceLevel = admissionapi.LevelPrivileged
		f2 := e2e.NewDefaultFramework("net-services2")
		// TODO(sur): verify if privileged is really necessary in a follow-up
		f2.NamespacePodSecurityEnforceLevel = admissionapi.LevelPrivileged

		It("should prevent connections to pods in different namespaces on the same node via service IPs", func() {
			Expect(checkServiceConnectivity(f1, f2, SAME_NODE)).NotTo(Succeed())
		})

		It("should prevent connections to pods in different namespaces on different nodes via service IPs", func() {
			Expect(checkServiceConnectivity(f1, f2, DIFFERENT_NODE)).NotTo(Succeed())
		})

		It("should allow connections to services in the default namespace from a pod in another namespace on the same node", func() {
			makeNamespaceGlobal(oc, f1.Namespace)
			Expect(checkServiceConnectivity(f1, f2, SAME_NODE)).To(Succeed())
		})
		It("should allow connections to services in the default namespace from a pod in another namespace on a different node", func() {
			makeNamespaceGlobal(oc, f1.Namespace)
			Expect(checkServiceConnectivity(f1, f2, DIFFERENT_NODE)).To(Succeed())
		})
		It("should allow connections from pods in the default namespace to a service in another namespace on the same node", func() {
			makeNamespaceGlobal(oc, f2.Namespace)
			Expect(checkServiceConnectivity(f1, f2, SAME_NODE)).To(Succeed())
		})
		It("should allow connections from pods in the default namespace to a service in another namespace on a different node", func() {
			makeNamespaceGlobal(oc, f2.Namespace)
			Expect(checkServiceConnectivity(f1, f2, DIFFERENT_NODE)).To(Succeed())
		})
	})
})
