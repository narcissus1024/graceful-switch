package data

import (
	"github.com/onsi/gomega"
	"testing"
)

func TestSSHConfigLoad(t *testing.T) {
	g := gomega.NewWithT(t)

	dm := GetDataManager()

	g.Expect(dm.InnerDataList.Load()).Should(gomega.BeNil())
	g.Expect(dm.InnerDataIndexList.Load()).Should(gomega.BeNil())
	g.Expect(dm.SSHConfig.Load()).Should(gomega.BeNil())
}

func TestMergeConfig(t *testing.T) {
	g := gomega.NewWithT(t)

	dm := GetDataManager()

	g.Expect(dm.InnerDataList.Load()).Should(gomega.BeNil())
	g.Expect(dm.SSHConfig.Load()).Should(gomega.BeNil())
	g.Expect(dm.MergeConfig()).Should(gomega.BeNil())
}