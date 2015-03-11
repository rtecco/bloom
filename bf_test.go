package bloom

import (
	"testing"

	"github.com/bmizerany/assert"
)

func TestBloom(t *testing.T) {

	bf := New(10, 0.05)

	bf.Add("stilton")
	bf.Add("colby")
	bf.Add("parmesean")

	assert.T(t, bf.Contains("parmesean"), "should contain 'parmesean'")
	assert.T(t, bf.Contains("stilton"), "should contain 'stilton''")
	assert.T(t, bf.Contains("colby"), "should contain 'colby'")
	assert.T(t, !bf.Contains("milk"), "shouldn't contain 'milk'")
}
