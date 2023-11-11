package sql

import (
	"github.com/aquasecurity/defsec/pkg/providers"
	"github.com/aquasecurity/defsec/pkg/providers/google/sql"
	"github.com/aquasecurity/defsec/pkg/scan"
	"github.com/aquasecurity/defsec/pkg/severity"
	"github.com/aquasecurity/defsec/pkg/state"
	"github.com/khulnasoft-lab/vul-policies/pkg/rules"
)

var CheckNoContainedDbAuth = rules.Register(
	scan.Rule{
		AVDID:       "AVD-GCP-0023",
		Provider:    providers.GoogleProvider,
		Service:     "sql",
		ShortCode:   "no-contained-db-auth",
		Summary:     "Contained database authentication should be disabled",
		Impact:      "Access can be granted without knowledge of the database administrator",
		Resolution:  "Disable contained database authentication",
		Explanation: `Users with ALTER permissions on users can grant access to a contained database without the knowledge of an administrator`,
		Links: []string{
			"https://docs.microsoft.com/en-us/sql/database-engine/configure-windows/contained-database-authentication-server-configuration-option?view=sql-server-ver15",
		},
		Terraform: &scan.EngineMetadata{
			GoodExamples:        terraformNoContainedDbAuthGoodExamples,
			BadExamples:         terraformNoContainedDbAuthBadExamples,
			Links:               terraformNoContainedDbAuthLinks,
			RemediationMarkdown: terraformNoContainedDbAuthRemediationMarkdown,
		},
		Severity: severity.Medium,
	},
	func(s *state.State) (results scan.Results) {
		for _, instance := range s.Google.SQL.Instances {
			if instance.Metadata.IsUnmanaged() {
				continue
			}
			if instance.DatabaseFamily() != sql.DatabaseFamilySQLServer {
				continue
			}
			if instance.Settings.Flags.ContainedDatabaseAuthentication.IsTrue() {
				results.Add(
					"Database instance has contained database authentication enabled.",
					instance.Settings.Flags.ContainedDatabaseAuthentication,
				)
			} else {
				results.AddPassed(&instance)
			}

		}
		return
	},
)
