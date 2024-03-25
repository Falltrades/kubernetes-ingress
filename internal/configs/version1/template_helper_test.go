package version1

import (
	"bytes"
	"strings"
	"testing"
	"text/template"
)

func TestMakeLocationPath_WithRegexCaseSensitiveModifier(t *testing.T) {
	t.Parallel()

	want := "~ \"^/coffee/[A-Z0-9]{3}\""
	got := makeLocationPath(
		&Location{Path: "/coffee/[A-Z0-9]{3}"},
		map[string]string{"nginx.org/path-regex": "case_sensitive"},
	)
	if got != want {
		t.Errorf("got: %s, want: %s", got, want)
	}
}

func TestMakeLocationPath_WithRegexCaseInsensitiveModifier(t *testing.T) {
	t.Parallel()

	want := "~* \"^/coffee/[A-Z0-9]{3}\""
	got := makeLocationPath(
		&Location{Path: "/coffee/[A-Z0-9]{3}"},
		map[string]string{"nginx.org/path-regex": "case_insensitive"},
	)
	if got != want {
		t.Errorf("got: %s, want: %s", got, want)
	}
}

func TestMakeLocationPath_WithRegexExactModifier(t *testing.T) {
	t.Parallel()

	want := "= \"/coffee\""
	got := makeLocationPath(
		&Location{Path: "/coffee"},
		map[string]string{"nginx.org/path-regex": "exact"},
	)
	if got != want {
		t.Errorf("got: %s, want: %s", got, want)
	}
}

func TestMakeLocationPath_WithBogusRegexModifier(t *testing.T) {
	t.Parallel()

	want := "/coffee"
	got := makeLocationPath(
		&Location{Path: "/coffee"},
		map[string]string{"nginx.org/path-regex": "bogus"},
	)
	if got != want {
		t.Errorf("got: %s, want: %s", got, want)
	}
}

func TestMakeLocationPath_WithEmptyRegexModifier(t *testing.T) {
	t.Parallel()

	want := "/coffee"
	got := makeLocationPath(
		&Location{Path: "/coffee"},
		map[string]string{"nginx.org/path-regex": ""},
	)
	if got != want {
		t.Errorf("got: %s, want: %s", got, want)
	}
}

func TestMakeLocationPath_WithBogusAnnotationName(t *testing.T) {
	t.Parallel()

	want := "/coffee"
	got := makeLocationPath(
		&Location{Path: "/coffee"},
		map[string]string{"nginx.org/bogus-annotation": ""},
	)
	if got != want {
		t.Errorf("got: %s, want: %s", got, want)
	}
}

func TestMakeLocationPath_ForIngressWithoutPathRegex(t *testing.T) {
	t.Parallel()

	want := "/coffee"
	got := makeLocationPath(
		&Location{Path: "/coffee"},
		map[string]string{},
	)
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestMakeLocationPath_ForIngressWithPathRegexCaseSensitive(t *testing.T) {
	t.Parallel()

	want := "~ \"^/coffee\""
	got := makeLocationPath(
		&Location{Path: "/coffee"},
		map[string]string{
			"nginx.org/path-regex": "case_sensitive",
		},
	)
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestMakeLocationPath_ForIngressWithPathRegexSetOnMinion(t *testing.T) {
	t.Parallel()

	want := "~ \"^/coffee\""
	got := makeLocationPath(
		&Location{
			Path: "/coffee",
			MinionIngress: &Ingress{
				Name:      "cafe-ingress-coffee-minion",
				Namespace: "default",
				Annotations: map[string]string{
					"nginx.org/mergeable-ingress-type": "minion",
					"nginx.org/path-regex":             "case_sensitive",
				},
			},
		},
		map[string]string{
			"nginx.org/mergeable-ingress-type": "master",
		},
	)

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestMakeLocationPath_ForIngressWithPathRegexSetOnMaster(t *testing.T) {
	t.Parallel()

	want := "~ \"^/coffee\""
	got := makeLocationPath(
		&Location{
			Path: "/coffee",
			MinionIngress: &Ingress{
				Name:      "cafe-ingress-coffee-minion",
				Namespace: "default",
			},
		},
		map[string]string{
			"nginx.org/mergeable-ingress-type": "master",
			"nginx.org/path-regex":             "case_sensitive",
		},
	)

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestMakeLocationPath_SetOnMinionTakesPrecedenceOverMaster(t *testing.T) {
	t.Parallel()

	want := "= \"/coffee\""
	got := makeLocationPath(
		&Location{
			Path: "/coffee",
			MinionIngress: &Ingress{
				Name:      "cafe-ingress-coffee-minion",
				Namespace: "default",
				Annotations: map[string]string{
					"nginx.org/mergeable-ingress-type": "minion",
					"nginx.org/path-regex":             "exact",
				},
			},
		},
		map[string]string{
			"nginx.org/mergeable-ingress-type": "master",
			"nginx.org/path-regex":             "case_sensitive",
		},
	)

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestMakeLocationPath_PathRegexSetOnMasterDoesNotModifyMinionWithoutPathRegexAnnotation(t *testing.T) {
	t.Parallel()

	want := "/coffee"
	got := makeLocationPath(
		&Location{
			Path: "/coffee",
			MinionIngress: &Ingress{
				Name:      "cafe-ingress-coffee-minion",
				Namespace: "default",
				Annotations: map[string]string{
					"nginx.org/mergeable-ingress-type": "minion",
				},
			},
		},
		map[string]string{
			"nginx.org/mergeable-ingress-type": "master",
			"nginx.org/path-regex":             "exact",
		},
	)

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestMakeLocationPath_ForIngress(t *testing.T) {
	t.Parallel()

	want := "~ \"^/coffee\""
	got := makeLocationPath(
		&Location{
			Path: "/coffee",
		},
		map[string]string{
			"nginx.org/path-regex": "case_sensitive",
		},
	)

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
func TestSplitInputString(t *testing.T) {
	t.Parallel()

	tmpl := newSplitTemplate(t)
	var buf bytes.Buffer

	input := "foo,bar"
	expected := "foo bar "

	err := tmpl.Execute(&buf, input)
	if err != nil {
		t.Fatalf("Failed to execute the template %v", err)
	}
	if buf.String() != expected {
		t.Errorf("Template generated wrong config, got %v but expected %v.", buf.String(), expected)
	}
}

func TestTrimWhiteSpaceFromInputString(t *testing.T) {
	t.Parallel()

	tmpl := newTrimTemplate(t)
	inputs := []string{
		"  foobar     ",
		"foobar   ",
		"   foobar",
		"foobar",
	}
	expected := "foobar"

	for _, i := range inputs {
		var buf bytes.Buffer
		err := tmpl.Execute(&buf, i)
		if err != nil {
			t.Fatalf("Failed to execute the template %v", err)
		}
		if buf.String() != expected {
			t.Errorf("Template generated wrong config, got %v but expected %v.", buf.String(), expected)
		}
	}
}

func TestContainsSubstring(t *testing.T) {
	t.Parallel()

	tmpl := newContainsTemplate(t)
	testCases := []struct {
		InputString string
		Substring   string
		expected    string
	}{
		{InputString: "foo", Substring: "foo", expected: "true"},
		{InputString: "foobar", Substring: "foo", expected: "true"},
		{InputString: "foo", Substring: "", expected: "true"},
		{InputString: "foo", Substring: "bar", expected: "false"},
		{InputString: "foo", Substring: "foobar", expected: "false"},
		{InputString: "", Substring: "foo", expected: "false"},
	}

	for _, tc := range testCases {
		var buf bytes.Buffer
		err := tmpl.Execute(&buf, tc)
		if err != nil {
			t.Fatalf("Failed to execute the template %v", err)
		}
		if buf.String() != tc.expected {
			t.Errorf("Template generated wrong config, got %v but expected %v.", buf.String(), tc.expected)
		}
	}
}

func TestHasPrefix(t *testing.T) {
	t.Parallel()

	tmpl := newHasPrefixTemplate(t)
	testCases := []struct {
		InputString string
		Prefix      string
		expected    string
	}{
		{InputString: "foo", Prefix: "foo", expected: "true"},
		{InputString: "foo", Prefix: "f", expected: "true"},
		{InputString: "foo", Prefix: "", expected: "true"},
		{InputString: "foo", Prefix: "oo", expected: "false"},
		{InputString: "foo", Prefix: "bar", expected: "false"},
		{InputString: "foo", Prefix: "foobar", expected: "false"},
	}

	for _, tc := range testCases {
		var buf bytes.Buffer
		err := tmpl.Execute(&buf, tc)
		if err != nil {
			t.Fatalf("Failed to execute the template %v", err)
		}
		if buf.String() != tc.expected {
			t.Errorf("Template generated wrong config, got %v but expected %v.", buf.String(), tc.expected)
		}
	}
}

func TestHasSuffix(t *testing.T) {
	t.Parallel()

	tmpl := newHasSuffixTemplate(t)
	testCases := []struct {
		InputString string
		Suffix      string
		expected    string
	}{
		{InputString: "bar", Suffix: "bar", expected: "true"},
		{InputString: "bar", Suffix: "r", expected: "true"},
		{InputString: "bar", Suffix: "", expected: "true"},
		{InputString: "bar", Suffix: "ba", expected: "false"},
		{InputString: "bar", Suffix: "foo", expected: "false"},
		{InputString: "bar", Suffix: "foobar", expected: "false"},
	}

	for _, tc := range testCases {
		var buf bytes.Buffer
		err := tmpl.Execute(&buf, tc)
		if err != nil {
			t.Fatalf("Failed to execute the template %v", err)
		}
		if buf.String() != tc.expected {
			t.Errorf("Template generated wrong config, got %v but expected %v.", buf.String(), tc.expected)
		}
	}
}

func TestToLowerInputString(t *testing.T) {
	t.Parallel()

	tmpl := newToLowerTemplate(t)
	testCases := []struct {
		InputString string
		expected    string
	}{
		{InputString: "foobar", expected: "foobar"},
		{InputString: "FOOBAR", expected: "foobar"},
		{InputString: "fOoBaR", expected: "foobar"},
		{InputString: "", expected: ""},
	}

	for _, tc := range testCases {
		var buf bytes.Buffer
		err := tmpl.Execute(&buf, tc)
		if err != nil {
			t.Fatalf("Failed to execute the template %v", err)
		}
		if buf.String() != tc.expected {
			t.Errorf("Template generated wrong config, got %v but expected %v.", buf.String(), tc.expected)
		}
	}
}

func TestToUpperInputString(t *testing.T) {
	t.Parallel()

	tmpl := newToUpperTemplate(t)
	testCases := []struct {
		InputString string
		expected    string
	}{
		{InputString: "foobar", expected: "FOOBAR"},
		{InputString: "FOOBAR", expected: "FOOBAR"},
		{InputString: "fOoBaR", expected: "FOOBAR"},
		{InputString: "", expected: ""},
	}

	for _, tc := range testCases {
		var buf bytes.Buffer
		err := tmpl.Execute(&buf, tc)
		if err != nil {
			t.Fatalf("Failed to execute the template %v", err)
		}
		if buf.String() != tc.expected {
			t.Errorf("Template generated wrong config, got %v but expected %v.", buf.String(), tc.expected)
		}
	}
}

func newSplitTemplate(t *testing.T) *template.Template {
	t.Helper()
	tmpl, err := template.New("testTemplate").Funcs(helperFunctions).Parse(`{{range $n := split . ","}}{{$n}} {{end}}`)
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}
	return tmpl
}

func newTrimTemplate(t *testing.T) *template.Template {
	t.Helper()
	tmpl, err := template.New("testTemplate").Funcs(helperFunctions).Parse(`{{trim .}}`)
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}
	return tmpl
}

func newContainsTemplate(t *testing.T) *template.Template {
	t.Helper()
	tmpl, err := template.New("testTemplate").Funcs(helperFunctions).Parse(`{{contains .InputString .Substring}}`)
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}
	return tmpl
}

func newHasPrefixTemplate(t *testing.T) *template.Template {
	t.Helper()
	tmpl, err := template.New("testTemplate").Funcs(helperFunctions).Parse(`{{hasPrefix .InputString .Prefix}}`)
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}
	return tmpl
}

func newHasSuffixTemplate(t *testing.T) *template.Template {
	t.Helper()
	tmpl, err := template.New("testTemplate").Funcs(helperFunctions).Parse(`{{hasSuffix .InputString .Suffix}}`)
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}
	return tmpl
}

func newToLowerTemplate(t *testing.T) *template.Template {
	t.Helper()
	tmpl, err := template.New("testTemplate").Funcs(helperFunctions).Parse(`{{toLower .InputString}}`)
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}
	return tmpl
}

func newToUpperTemplate(t *testing.T) *template.Template {
	t.Helper()
	tmpl, err := template.New("testTemplate").Funcs(helperFunctions).Parse(`{{toUpper .InputString}}`)
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}
	return tmpl
}

func TestValidateGenerateProxySetHeadersForValidHeaders(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		annotations      map[string]string
		wantProxyHeaders []string
	}{
		{
			name: "One Header",
			annotations: map[string]string{
				"nginx.org/proxy-set-headers": "X-Forwarded-ABC",
			},
			wantProxyHeaders: []string{
				"proxy_set_header X-Forwarded-ABC $http_x_forwarded_abc;",
			},
		},
		{
			name: "Two Headers",
			annotations: map[string]string{
				"nginx.org/proxy-set-headers": "X-Forwarded-ABC,BVC",
			},
			wantProxyHeaders: []string{
				"proxy_set_header X-Forwarded-ABC $http_x_forwarded_abc;",
				"proxy_set_header BVC $http_bvc;",
			},
		},
		{
			name: "Two Headers with One Value",
			annotations: map[string]string{
				"nginx.org/proxy-set-headers": "X-Forwarded-ABC,BVC test",
			},
			wantProxyHeaders: []string{
				"proxy_set_header X-Forwarded-ABC $http_x_forwarded_abc;",
				`proxy_set_header BVC "test";`,
			},
		},
		{
			name: "Three Headers",
			annotations: map[string]string{
				"nginx.org/proxy-set-headers": "X-Forwarded-ABC,BVC,X-Forwarded-Test",
			},
			wantProxyHeaders: []string{
				"proxy_set_header X-Forwarded-ABC $http_x_forwarded_abc;",
				"proxy_set_header BVC $http_bvc;",
				"proxy_set_header X-Forwarded-Test $http_x_forwarded_test;",
			},
		},
		{
			name: "Three Headers with Two Value",
			annotations: map[string]string{
				"nginx.org/proxy-set-headers": "X-Forwarded-ABC abc,BVC bat,X-Forwarded-Test",
			},
			wantProxyHeaders: []string{
				`proxy_set_header X-Forwarded-ABC "abc";`,
				`proxy_set_header BVC "bat";`,
				"proxy_set_header X-Forwarded-Test $http_x_forwarded_test;",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			generatedConfig, err := generateProxySetHeaders(tc.annotations)
			if err != nil {
				t.Fatal(err)
			}
			if len(tc.wantProxyHeaders) != strings.Count(generatedConfig, "\n") {
				t.Fatalf("expected %d config lines, got %d", len(tc.wantProxyHeaders), strings.Count(generatedConfig, "\n"))
			}

			for _, line := range tc.wantProxyHeaders {
				if !strings.Contains(generatedConfig, line) {
					t.Errorf("expected line %q not found in generated config", line)
				}
			}
		})
	}
}
