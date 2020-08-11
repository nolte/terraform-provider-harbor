resource "harbor_retention_policy" "cleanup" {
    scope {
        ref = harbor_project.main.id
    }

    rule {
        template = "always"
        tag_selectors {
            decoration = "matches"
            pattern    = "master"
            extras     = jsonencode({
                untagged: false
            })
        }
        scope_selectors {
            repository {
                kind       = "doublestar"
                decoration = "repoMatches"
                pattern    = "**"
            }
        }
    }

    rule {
        disabled = true
        template = "latestPulledN"
        params = {
            "latestPulledN"      = 15
            "nDaysSinceLastPush" = 7
        }
        tag_selectors {
            kind       = "doublestar"
            decoration = "excludes"
            pattern    = "master"
            extras     = jsonencode({
                untagged: false
            })
        }
        scope_selectors {
            repository {
                kind       = "doublestar"
                decoration = "repoExcludes"
                pattern    = "nginx-*"
            }
        }
    }

    trigger {
        settings {
            cron = "0 0 0 * * *"
        }
    }
}
