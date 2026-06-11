#noinspection CucumberUndefinedStep
Feature: Scale generator API should return the requested scale

  Scenario: API returns the intervals in JSON for the requested scale
    When I request a scale for a Turkish saz
    Then I am provided with the just interval ratios for the saz

  Scenario: API returns a Scala file for the requested scale
    When I request a Scala file for a Turkish saz scale
    Then I am provided with a Scala file containing the just interval ratios for the saz