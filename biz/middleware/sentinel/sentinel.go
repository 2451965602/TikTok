package sentinel

import (
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/flow"
	"tiktok/pkg/constants"
)

func Init() {

	err := sentinel.InitDefault()
	if err != nil {
		panic(err)
	}

	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "default",
			Threshold:              constants.SentinelThreshold,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			StatIntervalInMs:       constants.SentinelStatIntervalInMs,
		},
	})
	if err != nil {
		panic(err)
	}

}
