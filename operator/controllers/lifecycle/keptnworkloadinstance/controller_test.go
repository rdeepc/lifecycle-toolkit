package keptnworkloadinstance

import (
	"context"
	"fmt"
	"strings"
	"testing"

	klcv1alpha2 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha2"
	apicommon "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha2/common"
	controllercommon "github.com/keptn/lifecycle-toolkit/operator/controllers/common"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/common/fake"
	controllererrors "github.com/keptn/lifecycle-toolkit/operator/controllers/errors"
	interfacesfake "github.com/keptn/lifecycle-toolkit/operator/controllers/lifecycle/interfaces/fake"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
	testrequire "github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	k8sfake "sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func TestKeptnWorkloadInstanceReconciler_reconcileDeployment_FailedReplicaSet(t *testing.T) {

	rep := int32(1)
	replicasetFail := makeReplicaSet("myrep", "default", &rep, 0)

	fakeClient := k8sfake.NewClientBuilder().WithObjects(replicasetFail).Build()

	err := klcv1alpha2.AddToScheme(fakeClient.Scheme())
	testrequire.Nil(t, err)

	workloadInstance := makeWorkloadInstanceWithRef(replicasetFail.ObjectMeta, "ReplicaSet")

	err = fakeClient.Create(context.TODO(), workloadInstance)
	require.Nil(t, err)

	r := &KeptnWorkloadInstanceReconciler{
		Client: fakeClient,
	}

	keptnState, err := r.reconcileDeployment(context.TODO(), workloadInstance)
	testrequire.Nil(t, err)
	testrequire.Equal(t, apicommon.StateProgressing, keptnState)
}

func TestKeptnWorkloadInstanceReconciler_reconcileDeployment_UnavailableReplicaSet(t *testing.T) {

	rep := int32(1)
	replicasetFail := makeReplicaSet("myrep", "default", &rep, 0)

	// do not put the ReplicaSet into the cluster
	fakeClient := k8sfake.NewClientBuilder().WithObjects().Build()

	err := klcv1alpha2.AddToScheme(fakeClient.Scheme())
	testrequire.Nil(t, err)

	workloadInstance := makeWorkloadInstanceWithRef(replicasetFail.ObjectMeta, "ReplicaSet")

	err = fakeClient.Create(context.TODO(), workloadInstance)
	require.Nil(t, err)

	r := &KeptnWorkloadInstanceReconciler{
		Client: fakeClient,
	}

	keptnState, err := r.reconcileDeployment(context.TODO(), workloadInstance)
	testrequire.NotNil(t, err)
	testrequire.Equal(t, apicommon.StateUnknown, keptnState)
}

func TestKeptnWorkloadInstanceReconciler_reconcileDeployment_FailedStatefulSet(t *testing.T) {

	rep := int32(1)
	statefulsetFail := makeStatefulSet("mystat", "default", &rep, 0)

	fakeClient := k8sfake.NewClientBuilder().WithObjects(statefulsetFail).Build()

	err := klcv1alpha2.AddToScheme(fakeClient.Scheme())
	testrequire.Nil(t, err)

	workloadInstance := makeWorkloadInstanceWithRef(statefulsetFail.ObjectMeta, "StatefulSet")

	err = fakeClient.Create(context.TODO(), workloadInstance)
	require.Nil(t, err)

	r := &KeptnWorkloadInstanceReconciler{
		Client: fakeClient,
	}

	keptnState, err := r.reconcileDeployment(context.TODO(), workloadInstance)
	testrequire.Nil(t, err)
	testrequire.Equal(t, apicommon.StateProgressing, keptnState)
}

func TestKeptnWorkloadInstanceReconciler_reconcileDeployment_UnavailableStatefulSet(t *testing.T) {

	rep := int32(1)
	statefulSetFail := makeStatefulSet("mystat", "default", &rep, 0)

	// do not put the StatefulSet into the cluster
	fakeClient := k8sfake.NewClientBuilder().WithObjects().Build()

	err := klcv1alpha2.AddToScheme(fakeClient.Scheme())
	testrequire.Nil(t, err)

	workloadInstance := makeWorkloadInstanceWithRef(statefulSetFail.ObjectMeta, "StatefulSet")

	err = fakeClient.Create(context.TODO(), workloadInstance)
	require.Nil(t, err)

	r := &KeptnWorkloadInstanceReconciler{
		Client: fakeClient,
	}

	keptnState, err := r.reconcileDeployment(context.TODO(), workloadInstance)
	testrequire.NotNil(t, err)
	testrequire.Equal(t, apicommon.StateUnknown, keptnState)
}

func TestKeptnWorkloadInstanceReconciler_reconcileDeployment_FailedDaemonSet(t *testing.T) {

	daemonSetFail := makeDaemonSet("mystat", "default", 1, 0)

	fakeClient := k8sfake.NewClientBuilder().WithObjects(daemonSetFail).Build()

	err := klcv1alpha2.AddToScheme(fakeClient.Scheme())
	testrequire.Nil(t, err)

	workloadInstance := makeWorkloadInstanceWithRef(daemonSetFail.ObjectMeta, "DaemonSet")

	err = fakeClient.Create(context.TODO(), workloadInstance)
	require.Nil(t, err)

	r := &KeptnWorkloadInstanceReconciler{
		Client: fakeClient,
	}

	keptnState, err := r.reconcileDeployment(context.TODO(), workloadInstance)
	testrequire.Nil(t, err)
	testrequire.Equal(t, apicommon.StateProgressing, keptnState)
}

func TestKeptnWorkloadInstanceReconciler_reconcileDeployment_UnavailableDaemonSet(t *testing.T) {
	daemonSetFail := makeDaemonSet("mystat", "default", 1, 0)

	// do not put the DaemonSet into the cluster
	fakeClient := k8sfake.NewClientBuilder().WithObjects().Build()

	err := klcv1alpha2.AddToScheme(fakeClient.Scheme())
	testrequire.Nil(t, err)

	workloadInstance := makeWorkloadInstanceWithRef(daemonSetFail.ObjectMeta, "DaemonSet")

	err = fakeClient.Create(context.TODO(), workloadInstance)
	require.Nil(t, err)

	r := &KeptnWorkloadInstanceReconciler{
		Client: fakeClient,
	}

	keptnState, err := r.reconcileDeployment(context.TODO(), workloadInstance)
	testrequire.NotNil(t, err)
	testrequire.Equal(t, apicommon.StateUnknown, keptnState)
}

func TestKeptnWorkloadInstanceReconciler_reconcileDeployment_ReadyReplicaSet(t *testing.T) {

	rep := int32(1)
	replicaSet := makeReplicaSet("myrep", "default", &rep, 1)

	fakeClient := k8sfake.NewClientBuilder().WithObjects(replicaSet).Build()

	err := klcv1alpha2.AddToScheme(fakeClient.Scheme())
	testrequire.Nil(t, err)

	workloadInstance := makeWorkloadInstanceWithRef(replicaSet.ObjectMeta, "ReplicaSet")

	err = fakeClient.Create(context.TODO(), workloadInstance)
	require.Nil(t, err)

	r := &KeptnWorkloadInstanceReconciler{
		Client: fakeClient,
	}

	keptnState, err := r.reconcileDeployment(context.TODO(), workloadInstance)
	testrequire.Nil(t, err)
	testrequire.Equal(t, apicommon.StateSucceeded, keptnState)
}

func TestKeptnWorkloadInstanceReconciler_reconcileDeployment_ReadyStatefulSet(t *testing.T) {

	rep := int32(1)
	statefulSet := makeStatefulSet("mystat", "default", &rep, 1)

	fakeClient := k8sfake.NewClientBuilder().WithObjects(statefulSet).Build()

	err := klcv1alpha2.AddToScheme(fakeClient.Scheme())
	testrequire.Nil(t, err)

	workloadInstance := makeWorkloadInstanceWithRef(statefulSet.ObjectMeta, "StatefulSet")

	err = fakeClient.Create(context.TODO(), workloadInstance)
	require.Nil(t, err)

	r := &KeptnWorkloadInstanceReconciler{
		Client: fakeClient,
	}

	keptnState, err := r.reconcileDeployment(context.TODO(), workloadInstance)
	testrequire.Nil(t, err)
	testrequire.Equal(t, apicommon.StateSucceeded, keptnState)
}

func TestKeptnWorkloadInstanceReconciler_reconcileDeployment_ReadyDaemonSet(t *testing.T) {

	daemonSet := makeDaemonSet("mystat", "default", 1, 1)

	fakeClient := k8sfake.NewClientBuilder().WithObjects(daemonSet).Build()

	err := klcv1alpha2.AddToScheme(fakeClient.Scheme())
	testrequire.Nil(t, err)

	workloadInstance := makeWorkloadInstanceWithRef(daemonSet.ObjectMeta, "DaemonSet")

	err = fakeClient.Create(context.TODO(), workloadInstance)
	require.Nil(t, err)

	r := &KeptnWorkloadInstanceReconciler{
		Client: fakeClient,
	}

	keptnState, err := r.reconcileDeployment(context.TODO(), workloadInstance)
	testrequire.Nil(t, err)
	testrequire.Equal(t, apicommon.StateSucceeded, keptnState)
}

func TestKeptnWorkloadInstanceReconciler_reconcileDeployment_UnsupportedReferenceKind(t *testing.T) {

	fakeClient := k8sfake.NewClientBuilder().WithObjects().Build()

	err := klcv1alpha2.AddToScheme(fakeClient.Scheme())
	testrequire.Nil(t, err)

	workloadInstance := makeWorkloadInstanceWithRef(metav1.ObjectMeta{}, "Unknown")

	err = fakeClient.Create(context.TODO(), workloadInstance)
	require.Nil(t, err)

	r := &KeptnWorkloadInstanceReconciler{
		Client: fakeClient,
	}

	keptnState, err := r.reconcileDeployment(context.TODO(), workloadInstance)
	testrequire.ErrorIs(t, err, controllererrors.ErrUnsupportedWorkloadInstanceResourceReference)
	testrequire.Equal(t, apicommon.StateUnknown, keptnState)
}

func TestKeptnWorkloadInstanceReconciler_IsPodRunning(t *testing.T) {
	p1 := makeNominatedPod("pod1", "node1", v1.PodRunning)
	p2 := makeNominatedPod("pod2", "node1", v1.PodPending)
	podList := &v1.PodList{Items: []v1.Pod{p1, p2}}
	podList2 := &v1.PodList{Items: []v1.Pod{p2}}
	r := &KeptnWorkloadInstanceReconciler{
		Client: k8sfake.NewClientBuilder().WithLists(podList).Build(),
	}
	isPodRunning, err := r.isPodRunning(context.TODO(), klcv1alpha2.ResourceReference{UID: "pod1"}, "node1")
	testrequire.Nil(t, err)
	if !isPodRunning {
		t.Errorf("Wrong!")
	}

	r2 := &KeptnWorkloadInstanceReconciler{
		Client: k8sfake.NewClientBuilder().WithLists(podList2).Build(),
	}
	isPodRunning, err = r2.isPodRunning(context.TODO(), klcv1alpha2.ResourceReference{UID: "pod1"}, "node1")
	testrequire.Nil(t, err)
	if isPodRunning {
		t.Errorf("Wrong!")
	}

}

func makeNominatedPod(podName string, nodeName string, phase v1.PodPhase) v1.Pod {
	return v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: nodeName,
			Name:      podName,
			UID:       types.UID(podName),
		},
		Status: v1.PodStatus{
			Phase:             phase,
			NominatedNodeName: nodeName,
		},
	}
}

func makeReplicaSet(name string, namespace string, wanted *int32, available int32) *appsv1.ReplicaSet {

	return &appsv1.ReplicaSet{
		TypeMeta: metav1.TypeMeta{
			Kind: "ReplicaSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			UID:       types.UID(name),
		},
		Spec: appsv1.ReplicaSetSpec{
			Replicas: wanted,
		},
		Status: appsv1.ReplicaSetStatus{
			AvailableReplicas: available,
		},
	}

}

func makeStatefulSet(name string, namespace string, wanted *int32, available int32) *appsv1.StatefulSet {

	return &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			Kind: "StatefulSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			UID:       types.UID(name),
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: wanted,
		},
		Status: appsv1.StatefulSetStatus{
			AvailableReplicas: available,
		},
	}

}

func makeDaemonSet(name string, namespace string, wanted int32, available int32) *appsv1.DaemonSet {

	return &appsv1.DaemonSet{
		TypeMeta: metav1.TypeMeta{
			Kind: "StatefulSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			UID:       types.UID(name),
		},
		Spec: appsv1.DaemonSetSpec{},
		Status: appsv1.DaemonSetStatus{
			DesiredNumberScheduled: wanted,
			NumberReady:            available,
		},
	}

}

func Test_getAppVersionForWorkloadInstance(t *testing.T) {
	tests := []struct {
		name           string
		wli            *klcv1alpha2.KeptnWorkloadInstance
		list           *klcv1alpha2.KeptnAppVersionList
		wantFound      bool
		wantAppVersion klcv1alpha2.KeptnAppVersion
		wantErr        bool
	}{
		{
			name: "no appVersions",
			wli: &klcv1alpha2.KeptnWorkloadInstance{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-workloadinstance",
					Namespace: "default",
				},
				Spec: klcv1alpha2.KeptnWorkloadInstanceSpec{
					KeptnWorkloadSpec: klcv1alpha2.KeptnWorkloadSpec{
						AppName: "my-app",
						Version: "1.0",
					},
					WorkloadName: "my-app-my-workload",
				},
			},
			list:           &klcv1alpha2.KeptnAppVersionList{},
			wantFound:      false,
			wantAppVersion: klcv1alpha2.KeptnAppVersion{},
			wantErr:        false,
		},
		{
			name: "appVersion found",
			wli: &klcv1alpha2.KeptnWorkloadInstance{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-workloadinstance",
					Namespace: "default",
				},
				Spec: klcv1alpha2.KeptnWorkloadInstanceSpec{
					KeptnWorkloadSpec: klcv1alpha2.KeptnWorkloadSpec{
						AppName: "my-app",
						Version: "1.0",
					},
					WorkloadName: "my-app-my-workload",
				},
			},
			list: &klcv1alpha2.KeptnAppVersionList{
				Items: []klcv1alpha2.KeptnAppVersion{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "my-app",
							Namespace: "default",
						},
						Spec: klcv1alpha2.KeptnAppVersionSpec{
							KeptnAppSpec: klcv1alpha2.KeptnAppSpec{
								Version: "1.0",
								Workloads: []klcv1alpha2.KeptnWorkloadRef{
									{
										Name:    "my-workload",
										Version: "1.0",
									},
								},
							},
							AppName: "my-app",
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "my-app2",
							Namespace: "default",
						},
						Spec: klcv1alpha2.KeptnAppVersionSpec{
							KeptnAppSpec: klcv1alpha2.KeptnAppSpec{
								Version: "2.0",
								Workloads: []klcv1alpha2.KeptnWorkloadRef{
									{
										Name:    "my-workload",
										Version: "1.0",
									},
								},
							},
							AppName: "my-app",
						},
					},
				},
			},
			wantFound: true,
			wantAppVersion: klcv1alpha2.KeptnAppVersion{
				ObjectMeta: metav1.ObjectMeta{
					Name:            "my-app2",
					Namespace:       "default",
					ResourceVersion: "999",
				},
				Spec: klcv1alpha2.KeptnAppVersionSpec{
					KeptnAppSpec: klcv1alpha2.KeptnAppSpec{
						Version: "2.0",
						Workloads: []klcv1alpha2.KeptnWorkloadRef{
							{
								Name:    "my-workload",
								Version: "1.0",
							},
						},
					},
					AppName: "my-app",
				},
			},
			wantErr: false,
		},
		{
			name: "appVersion deprecated",
			wli: &klcv1alpha2.KeptnWorkloadInstance{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-workloadinstance",
					Namespace: "default",
				},
				Spec: klcv1alpha2.KeptnWorkloadInstanceSpec{
					KeptnWorkloadSpec: klcv1alpha2.KeptnWorkloadSpec{
						AppName: "my-app",
						Version: "1.0",
					},
					WorkloadName: "my-app-my-workload",
				},
			},
			list: &klcv1alpha2.KeptnAppVersionList{
				Items: []klcv1alpha2.KeptnAppVersion{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "my-app",
							Namespace: "default",
						},
						Spec: klcv1alpha2.KeptnAppVersionSpec{
							KeptnAppSpec: klcv1alpha2.KeptnAppSpec{
								Version: "1.0",
								Workloads: []klcv1alpha2.KeptnWorkloadRef{
									{
										Name:    "my-workload",
										Version: "1.0",
									},
								},
							},
							AppName: "my-app",
						},
						Status: klcv1alpha2.KeptnAppVersionStatus{
							Status: apicommon.StateDeprecated,
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "my-app2",
							Namespace: "default",
						},
						Spec: klcv1alpha2.KeptnAppVersionSpec{
							KeptnAppSpec: klcv1alpha2.KeptnAppSpec{
								Version: "2.0",
								Workloads: []klcv1alpha2.KeptnWorkloadRef{
									{
										Name:    "my-workload",
										Version: "1.0",
									},
								},
							},
							AppName: "my-app",
						},
						Status: klcv1alpha2.KeptnAppVersionStatus{
							Status: apicommon.StateDeprecated,
						},
					},
				},
			},
			wantFound:      false,
			wantAppVersion: klcv1alpha2.KeptnAppVersion{},
			wantErr:        false,
		},
		{
			name: "no workload for appversion",
			wli: &klcv1alpha2.KeptnWorkloadInstance{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-workloadinstance3",
					Namespace: "default",
				},
				Spec: klcv1alpha2.KeptnWorkloadInstanceSpec{
					KeptnWorkloadSpec: klcv1alpha2.KeptnWorkloadSpec{
						AppName: "my-app333",
						Version: "1.0.0",
					},
					WorkloadName: "my-app-my-workload",
				},
			},
			list: &klcv1alpha2.KeptnAppVersionList{
				Items: []klcv1alpha2.KeptnAppVersion{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "my-app",
							Namespace: "default",
						},
						Spec: klcv1alpha2.KeptnAppVersionSpec{
							KeptnAppSpec: klcv1alpha2.KeptnAppSpec{
								Version: "1.0",
								Workloads: []klcv1alpha2.KeptnWorkloadRef{
									{
										Name:    "my-workload",
										Version: "1.0",
									},
								},
							},
							AppName: "my-app",
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "my-app2",
							Namespace: "default",
						},
						Spec: klcv1alpha2.KeptnAppVersionSpec{
							KeptnAppSpec: klcv1alpha2.KeptnAppSpec{
								Version: "2.0",
								Workloads: []klcv1alpha2.KeptnWorkloadRef{
									{
										Name:    "my-workload",
										Version: "1.0",
									},
								},
							},
							AppName: "my-app",
						},
					},
				},
			},
			wantFound:      false,
			wantAppVersion: klcv1alpha2.KeptnAppVersion{},
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := klcv1alpha2.AddToScheme(scheme.Scheme)
			require.Nil(t, err)
			r := &KeptnWorkloadInstanceReconciler{
				Client: k8sfake.NewClientBuilder().WithLists(tt.list).Build(),
			}
			found, gotAppVersion, err := r.getAppVersionForWorkloadInstance(context.TODO(), tt.wli)
			require.Equal(t, tt.wantErr, err != nil)
			require.Equal(t, tt.wantFound, found)
			require.Equal(t, tt.wantAppVersion, gotAppVersion)
		})
	}
}

func Test_getLatestAppVersion(t *testing.T) {
	type args struct {
		apps *klcv1alpha2.KeptnAppVersionList
		wli  *klcv1alpha2.KeptnWorkloadInstance
	}
	tests := []struct {
		name           string
		args           args
		wantFound      bool
		wantAppVersion klcv1alpha2.KeptnAppVersion
		wantErr        bool
	}{
		{
			name: "app version found",
			args: args{
				apps: &klcv1alpha2.KeptnAppVersionList{
					Items: []klcv1alpha2.KeptnAppVersion{
						{
							ObjectMeta: metav1.ObjectMeta{
								Name:      "my-app",
								Namespace: "default",
							},
							Spec: klcv1alpha2.KeptnAppVersionSpec{
								KeptnAppSpec: klcv1alpha2.KeptnAppSpec{
									Version: "1.0",
									Workloads: []klcv1alpha2.KeptnWorkloadRef{
										{
											Name:    "my-workload",
											Version: "1.0",
										},
									},
								},
								AppName: "my-app",
							},
						},
						{
							ObjectMeta: metav1.ObjectMeta{
								Name:      "my-app",
								Namespace: "default",
							},
							Spec: klcv1alpha2.KeptnAppVersionSpec{
								KeptnAppSpec: klcv1alpha2.KeptnAppSpec{
									Version: "2.0",
									Workloads: []klcv1alpha2.KeptnWorkloadRef{
										{
											Name:    "my-workload",
											Version: "1.0",
										},
									},
								},
								AppName: "my-app",
							},
						},
					},
				},
				wli: &klcv1alpha2.KeptnWorkloadInstance{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "my-workloadinstance",
						Namespace: "default",
					},
					Spec: klcv1alpha2.KeptnWorkloadInstanceSpec{
						KeptnWorkloadSpec: klcv1alpha2.KeptnWorkloadSpec{
							AppName: "my-app",
							Version: "1.0",
						},
						WorkloadName: "my-app-my-workload",
					},
				},
			},
			wantFound: true,
			wantAppVersion: klcv1alpha2.KeptnAppVersion{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-app",
					Namespace: "default",
				},
				Spec: klcv1alpha2.KeptnAppVersionSpec{
					KeptnAppSpec: klcv1alpha2.KeptnAppSpec{
						Version: "2.0",
						Workloads: []klcv1alpha2.KeptnWorkloadRef{
							{
								Name:    "my-workload",
								Version: "1.0",
							},
						},
					},
					AppName: "my-app",
				},
			},
			wantErr: false,
		},
		{
			name: "app version not found",
			args: args{
				apps: &klcv1alpha2.KeptnAppVersionList{
					Items: []klcv1alpha2.KeptnAppVersion{
						{
							ObjectMeta: metav1.ObjectMeta{
								Name:      "my-app",
								Namespace: "default",
							},
							Spec: klcv1alpha2.KeptnAppVersionSpec{
								KeptnAppSpec: klcv1alpha2.KeptnAppSpec{
									Version: "1.0",
									Workloads: []klcv1alpha2.KeptnWorkloadRef{
										{
											Name:    "my-other-workload",
											Version: "1.0",
										},
									},
								},
								AppName: "my-app",
							},
						},
					},
				},
				wli: &klcv1alpha2.KeptnWorkloadInstance{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "my-workloadinstance",
						Namespace: "default",
					},
					Spec: klcv1alpha2.KeptnWorkloadInstanceSpec{
						KeptnWorkloadSpec: klcv1alpha2.KeptnWorkloadSpec{
							AppName: "my-app",
							Version: "1.0",
						},
						WorkloadName: "my-app-my-workload",
					},
				},
			},
			wantFound:      false,
			wantAppVersion: klcv1alpha2.KeptnAppVersion{},
			wantErr:        false,
		},
		{
			name: "app version with invalid version",
			args: args{
				apps: &klcv1alpha2.KeptnAppVersionList{
					Items: []klcv1alpha2.KeptnAppVersion{
						{
							ObjectMeta: metav1.ObjectMeta{
								Name:      "my-app",
								Namespace: "default",
							},
							Spec: klcv1alpha2.KeptnAppVersionSpec{
								KeptnAppSpec: klcv1alpha2.KeptnAppSpec{
									Version: "",
									Workloads: []klcv1alpha2.KeptnWorkloadRef{
										{
											Name:    "my-workload",
											Version: "1.0",
										},
									},
								},
								AppName: "my-app",
							},
						},
					},
				},
				wli: &klcv1alpha2.KeptnWorkloadInstance{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "my-workloadinstance",
						Namespace: "default",
					},
					Spec: klcv1alpha2.KeptnWorkloadInstanceSpec{
						KeptnWorkloadSpec: klcv1alpha2.KeptnWorkloadSpec{
							AppName: "my-app",
							Version: "1.0",
						},
						WorkloadName: "my-app-my-workload",
					},
				},
			},
			wantFound:      false,
			wantAppVersion: klcv1alpha2.KeptnAppVersion{},
			wantErr:        true,
		},
		{
			name: "app version list empty",
			args: args{
				apps: &klcv1alpha2.KeptnAppVersionList{
					Items: []klcv1alpha2.KeptnAppVersion{},
				},
				wli: &klcv1alpha2.KeptnWorkloadInstance{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "my-workloadinstance",
						Namespace: "default",
					},
					Spec: klcv1alpha2.KeptnWorkloadInstanceSpec{
						KeptnWorkloadSpec: klcv1alpha2.KeptnWorkloadSpec{
							AppName: "my-app",
							Version: "1.0",
						},
						WorkloadName: "my-app-my-workload",
					},
				},
			},
			wantFound:      false,
			wantAppVersion: klcv1alpha2.KeptnAppVersion{},
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found, gotAppVersion, err := getLatestAppVersion(tt.args.apps, tt.args.wli)
			require.Equal(t, tt.wantErr, err != nil)
			require.Equal(t, tt.wantFound, found)
			require.Equal(t, tt.wantAppVersion, gotAppVersion)
		})
	}
}

func TestKeptnWorkloadInstanceReconciler_ReconcileNoActionRequired(t *testing.T) {
	r, _, _ := setupReconciler()

	result, err := r.Reconcile(context.TODO(), ctrl.Request{})

	require.Nil(t, err)
	require.NotNil(t, result)
}

func TestKeptnWorkloadInstanceReconciler_ReconcileDoNotStartBeforeAppPreEvaluationIsDone(t *testing.T) {
	r, eventChannel, _ := setupReconciler()

	testNamespace := "some-ns"

	wi := &klcv1alpha2.KeptnWorkloadInstance{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "some-wi",
			Namespace: testNamespace,
		},
		Spec: klcv1alpha2.KeptnWorkloadInstanceSpec{
			KeptnWorkloadSpec: klcv1alpha2.KeptnWorkloadSpec{
				AppName: "some-app",
				Version: "1.0.0",
			},
			WorkloadName:    "some-app-some-workload",
			PreviousVersion: "",
			TraceId:         nil,
		},
		Status: klcv1alpha2.KeptnWorkloadInstanceStatus{},
	}

	err := r.Client.Create(context.TODO(), wi)

	require.Nil(t, err)

	err = controllercommon.AddAppVersion(
		r.Client,
		testNamespace,
		"some-app",
		"1.0.0",
		[]klcv1alpha2.KeptnWorkloadRef{
			{
				Name:    "some-workload",
				Version: "1.0.0",
			},
		},
		klcv1alpha2.KeptnAppVersionStatus{},
	)
	require.Nil(t, err)

	req := ctrl.Request{
		NamespacedName: types.NamespacedName{
			Namespace: testNamespace,
			Name:      "some-wi",
		},
	}

	result, err := r.Reconcile(context.TODO(), req)

	require.Nil(t, err)
	require.True(t, result.Requeue)

	expectedEvents := []string{
		"AppPreDeployEvaluationsNotFinished",
	}

	for _, e := range expectedEvents {
		event := <-eventChannel
		assert.Equal(t, strings.Contains(event, req.Name), true, "wrong appversion")
		assert.Equal(t, strings.Contains(event, req.Namespace), true, "wrong namespace")
		assert.Equal(t, strings.Contains(event, e), true, fmt.Sprintf("no %s found in %s", e, event))
	}
}

func TestKeptnWorkloadInstanceReconciler_ReconcileReachCompletion(t *testing.T) {
	r, eventChannel, _ := setupReconciler()

	testNamespace := "some-ns"

	wi := &klcv1alpha2.KeptnWorkloadInstance{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "some-wi",
			Namespace: testNamespace,
		},
		Spec: klcv1alpha2.KeptnWorkloadInstanceSpec{
			KeptnWorkloadSpec: klcv1alpha2.KeptnWorkloadSpec{
				AppName: "some-app",
				Version: "1.0.0",
			},
			WorkloadName:    "some-app-some-workload",
			PreviousVersion: "",
			TraceId:         nil,
		},
		Status: klcv1alpha2.KeptnWorkloadInstanceStatus{
			DeploymentStatus:               apicommon.StateSucceeded,
			PreDeploymentStatus:            apicommon.StateSucceeded,
			PostDeploymentStatus:           apicommon.StateSucceeded,
			PreDeploymentEvaluationStatus:  apicommon.StateSucceeded,
			PostDeploymentEvaluationStatus: apicommon.StateSucceeded,
			CurrentPhase:                   apicommon.PhaseWorkloadPostEvaluation.ShortName,
			Status:                         apicommon.StateSucceeded,
			StartTime:                      metav1.Time{},
			EndTime:                        metav1.Time{},
		},
	}

	err := r.Client.Create(context.TODO(), wi)

	require.Nil(t, err)

	err = controllercommon.AddAppVersion(
		r.Client,
		testNamespace,
		"some-app",
		"1.0.0",
		[]klcv1alpha2.KeptnWorkloadRef{
			{
				Name:    "some-workload",
				Version: "1.0.0",
			},
		},
		klcv1alpha2.KeptnAppVersionStatus{
			PreDeploymentEvaluationStatus: apicommon.StateSucceeded,
		},
	)
	require.Nil(t, err)

	req := ctrl.Request{
		NamespacedName: types.NamespacedName{
			Namespace: testNamespace,
			Name:      "some-wi",
		},
	}

	result, err := r.Reconcile(context.TODO(), req)

	require.Nil(t, err)

	// do not requeue since we reached completion
	require.False(t, result.Requeue)

	// here we do not expect an event about the application preEvaluation being finished since that  will have been sent in
	// one of the previous reconciliation loops that lead to the first phase being reached
	expectedEvents := []string{
		"WorkloadPostDeployEvaluationsFinished",
	}

	for _, e := range expectedEvents {
		event := <-eventChannel
		assert.Equal(t, strings.Contains(event, req.Name), true, "wrong appversion")
		assert.Equal(t, strings.Contains(event, req.Namespace), true, "wrong namespace")
		assert.Equal(t, strings.Contains(event, e), true, fmt.Sprintf("no %s found in %s", e, event))
	}
}

func setupReconciler() (*KeptnWorkloadInstanceReconciler, chan string, *interfacesfake.ITracerMock) {
	//setup logger
	opts := zap.Options{
		Development: true,
	}
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	//fake a tracer
	tr := &interfacesfake.ITracerMock{StartFunc: func(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
		return ctx, trace.SpanFromContext(ctx)
	}}

	fakeClient := fake.NewClient()
	recorder := record.NewFakeRecorder(100)
	r := &KeptnWorkloadInstanceReconciler{
		Client:      fakeClient,
		Scheme:      scheme.Scheme,
		Recorder:    recorder,
		Log:         ctrl.Log.WithName("test-appController"),
		Tracer:      tr,
		Meters:      controllercommon.InitAppMeters(),
		SpanHandler: &controllercommon.SpanHandler{},
	}
	return r, recorder.Events, tr
}

func makeWorkloadInstanceWithRef(objectMeta metav1.ObjectMeta, refKind string) *klcv1alpha2.KeptnWorkloadInstance {
	workloadInstance := &klcv1alpha2.KeptnWorkloadInstance{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-wli",
			Namespace: "default",
		},
		Spec: klcv1alpha2.KeptnWorkloadInstanceSpec{
			KeptnWorkloadSpec: klcv1alpha2.KeptnWorkloadSpec{
				ResourceReference: klcv1alpha2.ResourceReference{
					UID:  objectMeta.UID,
					Name: objectMeta.Name,
					Kind: refKind,
				},
			},
		},
	}
	return workloadInstance
}
