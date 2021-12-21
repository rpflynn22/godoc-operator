package main

import (
	"flag"
	"os"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/scheme"

	godocApi "github.com/rpflynn22/godoc-operator/internal/api/v1alpha1"
	"github.com/rpflynn22/godoc-operator/internal/controllers"
)

func main() {
	setupLog := ctrl.Log.WithName("setup")

	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	rtScheme := setupScheme()

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: rtScheme,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = (&controllers.RepoReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Repo")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

func setupScheme() *runtime.Scheme {
	rtScheme := runtime.NewScheme()
	schemeBuilder := &scheme.Builder{GroupVersion: godocApi.GroupVersion}
	schemeBuilder.Register(&godocApi.Repo{}, &godocApi.RepoList{})
	utilruntime.Must(clientgoscheme.AddToScheme(rtScheme))
	utilruntime.Must(schemeBuilder.AddToScheme(rtScheme))
	return rtScheme
}
