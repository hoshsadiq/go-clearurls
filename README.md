# go-clearurls
This is a go implementation of [ClearURLs](https://docs.clearurls.xyz). It uses the same json data that is used by the Firefox/Chrome addon.

## Example usage:
<!-- BEGIN EMBED FILE: example_test.go;go -->
```go
package clearurls_test

import (
	"fmt"

	"github.com/hoshsadiq/go-clearurls"
)

func ExampleURLCleaner_Clean() {
	cu := clearurls.New()

	cleanedUrl, err := cu.Clean("https://www.amazon.com/dp/exampleProduct/ref=sxin_0_pb?__mk_de_DE=ÅMÅŽÕÑ&keywords=tea&pd_rd_i=exampleProduct&pd_rd_r=8d39e4cd-1e4f-43db-b6e7-72e969a84aa5&pd_rd_w=1pcKM&pd_rd_wg=hYrNl&pf_rd_p=50bbfd25-5ef7-41a2-68d6-74d854b30e30&pf_rd_r=0GMWD0YYKA7XFGX55ADP&qid=1517757263&rnid=2914120011")
	if err != nil {
		panic(err)
	}

	fmt.Println(cleanedUrl)

	// Output: https://www.amazon.com/dp/exampleProduct
}
```
<!-- END EMBED FILE -->

# Licensing
Please note, while the go code is MIT licensed, the data set is licensed according to [GNU Lesser General Public License v3.0](https://github.com/ClearURLs/Rules/blob/master/LICENSE) as per the license in the original repository.
