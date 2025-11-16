package metrics

import "testing"

// TestTemplate ensures the template for the homepage of unbound_exporter renders and properly escapes
func TestTemplate(t *testing.T) {
	const expected = `
<!DOCTYPE html>
<html>
	<head><title>Unbound Exporter</title></head>
	<body>
		<h1>Unbound Exporter</h1>
		<p><a href='/metrics%27%20bad=%27true%27'>Metrics</a></p>
		<p><a href='/_healthz'>Health</a></p>
	</body>
</html>
`

	result := string(homePageText("/metrics' bad='true'", "/_healthz"))
	if result != expected {
		t.Fatalf("Unexpected result: '%s' is not '%s'", result, expected)
	}
}
