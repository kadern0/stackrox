package multipliers

import (
	"time"

	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/protoconv"
)

func getMockDeployment() *storage.Deployment {
	return &storage.Deployment{
		ClusterId: "cluster",
		Containers: []*storage.Container{
			{
				Volumes: []*storage.Volume{
					{
						Name:     "readonly",
						ReadOnly: true,
					},
				},
				Secrets: []*storage.EmbeddedSecret{
					{
						Name: "secret",
					},
				},
				SecurityContext: &storage.SecurityContext{
					AddCapabilities: []string{
						"ALL",
					},
					Privileged: true,
				},
				Image: &storage.Image{
					Name: &storage.ImageName{
						FullName: "docker.io/library/nginx:1.10",
						Registry: "docker.io",
						Remote:   "library/nginx",
						Tag:      "1.10",
					},
					Scan: &storage.ImageScan{
						Components: []*storage.ImageScanComponent{
							{
								Name:    "comp1",
								Version: "1.1.1",
								Vulns: []*storage.Vulnerability{
									{
										Cvss: 5,
									},
									{
										Cvss: 5,
									},
								},
							},
							{
								Name:    "comp1",
								Version: "1.1.1",
								Vulns: []*storage.Vulnerability{
									{
										Cvss: 5,
									},
									{
										Cvss: 5,
									},
								},
							},
						},
					},
					Metadata: &storage.ImageMetadata{
						V1: &storage.V1Metadata{
							Created: protoconv.ConvertTimeToTimestamp(time.Now().Add(-(180 * 24 * time.Hour))),
						},
					},
				},
				Ports: []*storage.PortConfig{
					{
						Name:          "Port1",
						ContainerPort: 22,
						Exposure:      storage.PortConfig_EXTERNAL,
						ExposedPort:   8082,
					},
					{
						Name:          "Port2",
						ContainerPort: 23,
						Exposure:      storage.PortConfig_INTERNAL,
						ExposedPort:   8083,
					},
					{
						Name:          "Port3",
						ContainerPort: 24,
						Exposure:      storage.PortConfig_NODE,
						ExposedPort:   8084,
					},
				},
			},
			{
				Image: &storage.Image{
					Metadata: &storage.ImageMetadata{
						V1: &storage.V1Metadata{
							Created: protoconv.ConvertTimeToTimestamp(time.Now().Add(-(90 * 24 * time.Hour))),
						},
					},
				},
				Volumes: []*storage.Volume{
					{
						Name: "rw volume",
					},
				},
				SecurityContext: &storage.SecurityContext{},
			},
		},
	}
}
