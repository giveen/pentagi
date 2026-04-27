import type { CodegenConfig } from '@graphql-codegen/cli';

const config: CodegenConfig = {
    documents: './graphql-schema.graphql',
    generates: {
        './src/graphql/types.ts': {
            config: {
                apolloReactCommonImportFrom: '@apollo/client/core',
                apolloReactHooksImportFrom: '@apollo/client/react',
                dedupeFragments: true,
                exportFragmentSpreadSubTypes: true,
                inlineFragmentTypes: 'combine',
                preResolveTypes: true,
                skipTypename: true,
                useTypeImports: true,
                withHooks: true,
            },
            plugins: ['typescript', 'typescript-operations', 'typescript-react-apollo'],
        },
    },
    hooks: {
        afterOneFileWrite: ['npx prettier --write'],
    },
    schema: '../backend/pkg/graph/schema.graphqls',
};

export default config;
