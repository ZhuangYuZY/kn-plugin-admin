// Copyright © 2019 The Knative Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package domain

import (
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1beta1 "k8s.io/apimachinery/pkg/apis/meta/v1beta1"

	hprinters "knative.dev/client/pkg/printers"
)

// DomainListHandlers adds print handlers for domain list command
func DomainListHandlers(h hprinters.PrintHandler) {
	kDomainColumnDefinitions := []metav1beta1.TableColumnDefinition{
		{Name: "Custom-Domain", Type: "string", Description: "Name of Knative custom domain.", Priority: 1},
		{Name: "Selector", Type: "string", Description: "Selector of the Knative custom domain.", Priority: 1},
	}
	h.TableHandler(kDomainColumnDefinitions, printKDomainList)
}

// printKDomainList populates the Knative custom domain list table rows
func printKDomainList(domainCM *corev1.ConfigMap, options hprinters.PrintOptions) ([]metav1beta1.TableRow, error) {
	kDomainList := domainCM.Data
	delete(kDomainList, "_example")
	rows := make([]metav1beta1.TableRow, 0, len(kDomainList))
	for k, v := range kDomainList {
		row := metav1beta1.TableRow{}
		row.Cells = append(row.Cells, k, formatSelectorForPrint(v))
		rows = append(rows, []metav1beta1.TableRow{row}...)
	}
	return rows, nil
}

//format change for selector from "selector:\n  key1: value1\n  key2: value2\n" to "key1=value1; key2=value2;
func formatSelectorForPrint(selector string) string {
	parts := strings.Split(strings.ReplaceAll(strings.TrimSpace(selector), ":", "="), "\n")
	selectorForPrint := ""
	for i, v := range parts {
		if i == 0 && strings.Compare(v, "selector=") != 0 {
			return ""
		} else if i > 0 {
			if strings.Contains(v, "=") {
				selectorForPrint = strings.Join([]string{selectorForPrint, strings.ReplaceAll(v, " ", "")}, "")
				selectorForPrint = strings.Join([]string{selectorForPrint, "; "}, "")
			}
		}

	}
	return selectorForPrint
}
