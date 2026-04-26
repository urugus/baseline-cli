#!/usr/bin/env ruby
# frozen_string_literal: true

require "pathname"

formula_path = Pathname(ARGV.fetch(0))
version = ARGV.fetch(1)
checksums_path = Pathname(ARGV.fetch(2))

unless version.match?(/\Av\d+\.\d+\.\d+\z/)
  abort "version must look like vX.Y.Z: #{version}"
end

checksums = checksums_path.read.lines.each_with_object({}) do |line, result|
  sha, file = line.split
  result[file] = sha
end

required = %w[
  baseline-darwin-amd64
  baseline-darwin-arm64
  baseline-linux-amd64
  baseline-linux-arm64
]

missing = required.reject { |file| checksums[file] }
abort "missing checksums: #{missing.join(", ")}" unless missing.empty?

formula = formula_path.read
formula.gsub!(/version "\d+\.\d+\.\d+"/, %(version "#{version.delete_prefix("v")}"))
formula.gsub!(/download\/v\d+\.\d+\.\d+\//, "download/#{version}/")

replacements = {
  "baseline-darwin-arm64" => checksums.fetch("baseline-darwin-arm64"),
  "baseline-darwin-amd64" => checksums.fetch("baseline-darwin-amd64"),
  "baseline-linux-arm64" => checksums.fetch("baseline-linux-arm64"),
  "baseline-linux-amd64" => checksums.fetch("baseline-linux-amd64")
}

replacements.each do |file, sha|
  pattern = /(url ".*#{Regexp.escape(file)}"\n\s+sha256 ")[0-9a-f]+(")/
  formula.gsub!(pattern, "\\1#{sha}\\2")
end

formula.gsub!(/baseline v\d+\.\d+\.\d+/, "baseline #{version}")
formula_path.write(formula)
