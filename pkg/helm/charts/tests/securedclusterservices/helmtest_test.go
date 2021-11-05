package securedclusterservices

import (
	"testing"

	helmTest "github.com/stackrox/helmtest/pkg/framework"
	"github.com/stackrox/rox/image"
	"github.com/stretchr/testify/require"
	"helm.sh/helm/v3/pkg/chartutil"
)

func TestWithHelmtest(t *testing.T) {
	helmImage := image.GetDefaultImage()
	tpl, err := helmImage.GetSecuredClusterServicesChartTemplate()
	require.NoError(t, err, "error retrieving chart template")
	ch, err := tpl.InstantiateAndLoad(metaValues)
	require.NoError(t, err, "error instantiating chart")

	suite, err := helmTest.LoadSuite("testdata/helmtest")
	require.NoError(t, err, "failed to load helmtest suite")

	target := &helmTest.Target{
		Chart: ch,
		ReleaseOptions: chartutil.ReleaseOptions{
			Name:      "stackrox-secured-cluster-services",
			Namespace: "stackrox",
			IsInstall: true,
		},
	}
	suite.Run(t, target)
}
