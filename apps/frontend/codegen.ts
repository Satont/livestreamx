import { join, resolve } from 'node:path'
import process from 'node:process'

import type { CodegenConfig } from '@graphql-codegen/cli'

const schemaDir = resolve(join(process.cwd(), '..', 'api', 'schema', '*.graphqls'))

const config: CodegenConfig = {
	config: {
		scalars: {
			Upload: 'File',
		},
	},
	schema: schemaDir,
	documents: ['src/**/*.{ts,vue}'],
	ignoreNoDocuments: true, // for better experience with the watcher
	generates: {
		'./src/gql/': {
			preset: 'client',
			config: {
				useTypeImports: true,
			},
		},
	},
}

export default config