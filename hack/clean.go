package main

import (
	"log"

	uhc "github.com/openshift-online/uhc-sdk-go/pkg/client"
	"github.com/openshift-online/uhc-sdk-go/pkg/client/clustersmgmt/v1"

	"github.com/openshift/osde2e/pkg/config"
)

var OSD *uhc.Connection

func init() {
	logger, err := uhc.NewGoLoggerBuilder().
		Build()
	if err != nil {
		panic(err)
	}

	builder := uhc.NewConnectionBuilder().
		Logger(logger).
		Tokens(config.Cfg.UHCToken)

	OSD, err = builder.Build()
	if err != nil {
		panic(err)
	}
}

func main() {
	resp, err := OSD.ClustersMgmt().V1().Clusters().List().Send()
	if err != nil {
		panic(err)
	}

	log.Printf("Found %d clusters", resp.Total())

	log.Println("Pruning....")
	resp.Items().Each(func(c *v1.Cluster) bool {
		if c.State() == v1.ClusterStateReady {
			log.Println("Pruning " + c.ID() + "....")
			resp, err := OSD.ClustersMgmt().V1().Clusters().Cluster(c.ID()).Delete().Send()
			if err != nil {
				panic(err)
			}

			if resp.Status() != 200 {
				log.Println(resp.Error())
			}
		}
		return true
	})

	resp, err = OSD.ClustersMgmt().V1().Clusters().List().Send()
	if err != nil {
		panic(err)
	}

	log.Printf("Found %d clusters", resp.Total())
}
