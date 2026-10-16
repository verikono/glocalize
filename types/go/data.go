package types

// TypeAliases in this domain

type Schema = string

type (
	// TypeAlias for a string containing a GraphQL Schema. @see https://spec.graphql.org/September2025/
	GraphQLSchema = Schema

	// TypeAlias for a string containing a PSL. @see https://www.prisma.io/docs/orm/v7/reference/prisma-schema-reference
	PrismaSchema = Schema
)

type SchemaParser = func(schema Schema) (*Schema, error)
