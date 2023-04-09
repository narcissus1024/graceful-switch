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

func TestIsHostConfigItem(t *testing.T) {
	g := gomega.NewWithT(t)
	dm := GetDataManager()

	b := false
	v := ""

	b, v = dm.IsHostConfigItem("host aaa")
	g.Expect(b).Should(gomega.BeTrue())
	g.Expect(v).Should(gomega.Equal("aaa"))

	b, v = dm.IsHostConfigItem("HOST aaa")
	g.Expect(b).Should(gomega.BeTrue())
	g.Expect(v).Should(gomega.Equal("aaa"))

	b, v = dm.IsHostConfigItem("Host aaa")
	g.Expect(b).Should(gomega.BeTrue())
	g.Expect(v).Should(gomega.Equal("aaa"))

	b, v = dm.IsHostConfigItem("hOsT aaa")
	g.Expect(b).Should(gomega.BeTrue())
	g.Expect(v).Should(gomega.Equal("aaa"))

	b, v = dm.IsHostConfigItem("host aaa  ")
	g.Expect(b).Should(gomega.BeTrue())
	g.Expect(v).Should(gomega.Equal("aaa"))

	b, v = dm.IsHostConfigItem("	host aaa")
	g.Expect(b).Should(gomega.BeTrue())
	g.Expect(v).Should(gomega.Equal("aaa"))

	b, v = dm.IsHostConfigItem(" host aaa")
	g.Expect(b).Should(gomega.BeTrue())
	g.Expect(v).Should(gomega.Equal("aaa"))

	b, v = dm.IsHostConfigItem("  	host aaa")
	g.Expect(b).Should(gomega.BeTrue())
	g.Expect(v).Should(gomega.Equal("aaa"))

	b, v = dm.IsHostConfigItem("    Host ")
	g.Expect(b).Should(gomega.BeTrue())
	g.Expect(v).Should(gomega.Equal(""))

	b, v = dm.IsHostConfigItem("  	host  	 aaa")
	g.Expect(b).Should(gomega.BeTrue())
	g.Expect(v).Should(gomega.Equal("aaa"))

	b, v = dm.IsHostConfigItem("  	host=aaa")
	g.Expect(b).Should(gomega.BeTrue())
	g.Expect(v).Should(gomega.Equal("aaa"))

	b, v = dm.IsHostConfigItem("  	host  =aaa")
	g.Expect(b).Should(gomega.BeTrue())
	g.Expect(v).Should(gomega.Equal("aaa"))

	b, v = dm.IsHostConfigItem("  	host=  aaa")
	g.Expect(b).Should(gomega.BeTrue())
	g.Expect(v).Should(gomega.Equal("aaa"))

	b, v = dm.IsHostConfigItem("  	host  =  	aaa")
	g.Expect(b).Should(gomega.BeTrue())
	g.Expect(v).Should(gomega.Equal("aaa"))

	b, v = dm.IsHostConfigItem("  	host  =  	aaa  host")
	g.Expect(b).Should(gomega.BeTrue())
	g.Expect(v).Should(gomega.Equal("aaa  host"))

	b, _ = dm.IsHostConfigItem("aHost aaa")
	g.Expect(b).Should(gomega.BeFalse())

	b, _ = dm.IsHostConfigItem("Host == aaa")
	g.Expect(b).Should(gomega.BeFalse())

}