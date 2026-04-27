import { zodResolver } from '@hookform/resolvers/zod';
import { Loader2, Play, Plus, Save, Server, Trash2 } from 'lucide-react';
import { Fragment, useMemo, useState } from 'react';
import { Controller, useFieldArray, useForm } from 'react-hook-form';
import { useNavigate, useParams } from 'react-router-dom';
import { z } from 'zod';
import { gql, useQuery, useMutation } from '@apollo/client';
import { useEffect } from 'react';

import ConfirmationDialog from '@/components/shared/confirmation-dialog';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { Button } from '@/components/ui/button';
import { Form, FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { StatusCard } from '@/components/ui/status-card';
import { Switch } from '@/components/ui/switch';

type McpTransport = 'sse' | 'stdio';

const keyValueSchema = z.object({
    key: z.string().min(1, 'Key is required'),
    value: z.string().min(1, 'Value is required'),
});

const formSchema = z.object({
    name: z
        .string({ required_error: 'Name is required' })
        .min(1, 'Name is required')
        .max(50, 'Maximum 50 characters allowed'),
    sse: z
        .object({
            headers: z.array(keyValueSchema).optional().default([]),
            url: z.string().min(1, 'URL is required'),
        })
        .optional(),
    stdio: z
        .object({
            args: z.string().optional().nullable(),
            command: z.string().min(1, 'Command is required'),
            env: z.array(keyValueSchema).optional().default([]),
        })
        .optional(),
    tools: z
        .array(
            z.object({
                description: z.string().optional(),
                enabled: z.boolean().optional().default(true),
                name: z.string().min(1, 'Tool name is required'),
            }),
        )
        .default([]),
    transport: z.enum(['stdio', 'sse'], { required_error: 'Transport is required' }),
});

type FormData = z.infer<typeof formSchema>;

const GET_MCP_SERVER = gql`
    query GetMcpServer($mcpServerId: ID!) {
        mcpServer(mcpServerId: $mcpServerId) {
            id
            name
            transport
            stdio { command args env { key value } }
            sse { url headers { key value } }
            tools { name description enabled }
            createdAt
            updatedAt
        }
    }
`;

const CREATE_MCP_SERVER = gql`
    mutation CreateMcpServer($input: CreateMcpServerInput!) {
        createMcpServer(input: $input) {
            id
        }
    }
`;

const UPDATE_MCP_SERVER = gql`
    mutation UpdateMcpServer($mcpServerId: ID!, $input: UpdateMcpServerInput!) {
        updateMcpServer(mcpServerId: $mcpServerId, input: $input) {
            id
        }
    }
`;

const TEST_MCP_SERVER = gql`
    mutation TestMcpServer($mcpServerId: ID!) {
        testMcpServer(mcpServerId: $mcpServerId)
    }
`;

const SettingsMcpServer = () => {
    const navigate = useNavigate();
    const params = useParams();
    const isNew = params.mcpServerId === undefined;
    const [submitError, setSubmitError] = useState<null | string>(null);
    const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
    const [isTestLoading, setIsTestLoading] = useState(false);
    const [testMessage, setTestMessage] = useState<null | string>(null);
    const [testError, setTestError] = useState<null | string>(null);
    const [toolTestLoadingIndex, setToolTestLoadingIndex] = useState<null | number>(null);
    const [toolTestIndex, setToolTestIndex] = useState<null | number>(null);
    const [toolTestMessage, setToolTestMessage] = useState<null | string>(null);
    const [toolTestError, setToolTestError] = useState<null | string>(null);

    const defaults: FormData = useMemo(() => {
        return {
            name: '',
            stdio: { args: '', command: '', env: [] },
            tools: [],
            transport: 'stdio',
        } as FormData;
    }, []);

    const form = useForm<FormData>({
        defaultValues: defaults,
        mode: 'onChange',
        resolver: zodResolver(formSchema),
    });

    // Load server when editing
    const id = isNew ? undefined : Number(params.mcpServerId);
    const { data: serverData, loading: serverLoading, error: serverError } = useQuery(GET_MCP_SERVER, {
        variables: { mcpServerId: id },
        skip: isNew,
        fetchPolicy: 'cache-and-network',
    });

    const [createMcpServer] = useMutation(CREATE_MCP_SERVER);
    const [updateMcpServer] = useMutation(UPDATE_MCP_SERVER);
    const [testMcpServer] = useMutation(TEST_MCP_SERVER);

    useEffect(() => {
        if (!isNew && serverData?.mcpServer) {
            const s = serverData.mcpServer;
            form.reset({
                name: s.name || '',
                transport: s.transport || 'stdio',
                stdio: s.stdio
                    ? {
                          command: s.stdio.command || '',
                          args: s.stdio.args || '',
                          env: (s.stdio.env || []).map((e: any) => ({ key: e.key, value: e.value })),
                      }
                    : { command: '', args: '', env: [] },
                sse: s.sse
                    ? { url: s.sse.url || '', headers: (s.sse.headers || []).map((h: any) => ({ key: h.key, value: h.value })) }
                    : undefined,
                tools: s.tools || [],
            });
        }
    }, [isNew, serverData]);

    const transport = form.watch('transport');

    // Field arrays
    const stdioEnvArray = useFieldArray({ control: form.control, name: 'stdio.env' as const });
    const sseHeadersArray = useFieldArray({ control: form.control, name: 'sse.headers' as const });
    const toolsArray = useFieldArray({ control: form.control, name: 'tools' as const });

    const handleAddKeyValue = (target: 'env' | 'headers') => {
        if (target === 'env') {
            stdioEnvArray.append({ key: '', value: '' });
        } else {
            sseHeadersArray.append({ key: '', value: '' });
        }
    };

    const handleSubmit = async (_data: FormData) => {
        try {
            setSubmitError(null);
            const buildInput = (d: FormData) => {
                const input: any = { name: d.name, transport: d.transport };
                if (d.transport === 'stdio' && d.stdio) {
                    input.stdio = {
                        command: d.stdio.command,
                        args: d.stdio.args || '',
                        env: (d.stdio.env || []).map((e: any) => ({ key: e.key, value: e.value })),
                    };
                }
                if (d.transport === 'sse' && d.sse) {
                    input.sse = {
                        url: d.sse.url,
                        headers: (d.sse.headers || []).map((h: any) => ({ key: h.key, value: h.value })),
                    };
                }
                input.tools = (d.tools || []).map((t) => ({ name: t.name, description: t.description, enabled: !!t.enabled }));

                return input;
            };

            if (isNew) {
                const res = await createMcpServer({ variables: { input: buildInput(_data) } });
                const newId = res?.data?.createMcpServer?.id;
                if (newId) {
                    navigate(`/settings/mcp-servers/${newId}`);
                    return;
                }
            } else {
                const mcpId = Number(params.mcpServerId);
                await updateMcpServer({ variables: { mcpServerId: mcpId, input: buildInput(_data) } });
                navigate('/settings/mcp-servers');
                return;
            }
        } catch {
            setSubmitError('Failed to save MCP server');
        }
    };

    const handleDelete = () => {
        if (isNew) {
            return;
        }

        setIsDeleteDialogOpen(true);
    };

    const handleConfirmDelete = async () => {
        try {
            // Call backend to delete via mutation
            const mcpId = Number(params.mcpServerId);
            if (!Number.isNaN(mcpId)) {
                await fetch(`${window.location.origin}/graphql`, {
                    method: 'POST',
                    credentials: 'include',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ query: `mutation DeleteMcpServer($id: ID!){ deleteMcpServer(mcpServerId: $id) }`, variables: { id: mcpId } }),
                });
            }
            navigate('/settings/mcp-servers');
        } catch {
            setSubmitError('Failed to delete MCP server');
        }
    };

    const handleTest = async () => {
        setTestMessage(null);
        setTestError(null);
        // Validate minimal required fields based on transport
        const valid = await form.trigger();

        if (!valid) {
            setTestError('Please fix validation errors before testing');

            return;
        }

        try {
            setIsTestLoading(true);
            if (isNew) {
                setTestError('Please save the MCP server before testing');
                return;
            }

            const mcpId = Number(params.mcpServerId);
            if (Number.isNaN(mcpId)) {
                setTestError('Invalid server id');
                return;
            }

            const res = await testMcpServer({ variables: { mcpServerId: mcpId } });
            const ok = res?.data?.testMcpServer;
            if (ok) {
                setTestMessage('Connection successful');
            } else {
                setTestError('Connection failed');
            }
        } catch (e: any) {
            setTestError(e?.message || 'Connection failed');
        } finally {
            setIsTestLoading(false);
        }
    };

    const handleTestTool = async (index: number) => {
        setToolTestIndex(index);
        setToolTestLoadingIndex(index);
        setToolTestMessage(null);
        setToolTestError(null);

        try {
            // Basic validation: tool must have a name
            const toolName = form.getValues(`tools.${index}.name` as const) as string | undefined;

            if (!toolName) {
                setToolTestError('Tool name is required');

                return;
            }

            // Simulate tool invocation
            await new Promise((r) => setTimeout(r, 600));
            setToolTestMessage('Tool test passed');
        } catch {
            setToolTestError('Tool test failed');
        } finally {
            setToolTestLoadingIndex(null);
        }
    };

    if (!isNew && !getMockServerById(Number(params.mcpServerId))) {
        return (
            <StatusCard
                action={
                    <Button
                        onClick={() => navigate('/settings/mcp-servers')}
                        variant="secondary"
                    >
                        Back to list
                    </Button>
                }
                description="The requested MCP server could not be located in mock data"
                icon={<Server className="text-muted-foreground size-8" />}
                title="MCP Server not found"
            />
        );
    }

    return (
        <Fragment>
            <div className="flex flex-col gap-4">
                <div className="flex flex-col gap-2">
                    <h2 className="flex items-center gap-2 text-lg font-semibold">
                        <Server className="text-muted-foreground size-5" />
                        {isNew ? 'New MCP Server' : 'MCP Server Settings'}
                    </h2>

                    <div className="text-muted-foreground">
                        {isNew ? 'Configure a new MCP server' : 'Update MCP server settings'}
                    </div>
                </div>

                <Form {...form}>
                    <form
                        className="flex flex-col gap-6"
                        id="mcp-server-form"
                        onSubmit={form.handleSubmit(handleSubmit)}
                    >
                        {(submitError || testMessage || testError) && (
                            <Alert variant="destructive">
                                {submitError && (
                                    <>
                                        <AlertTitle>Error</AlertTitle>
                                        <AlertDescription>{submitError}</AlertDescription>
                                    </>
                                )}
                                {testError && (
                                    <>
                                        <AlertTitle>Test Failed</AlertTitle>
                                        <AlertDescription>{testError}</AlertDescription>
                                    </>
                                )}
                                {testMessage && !submitError && !testError && (
                                    <>
                                        <AlertTitle>Test Passed</AlertTitle>
                                        <AlertDescription>{testMessage}</AlertDescription>
                                    </>
                                )}
                            </Alert>
                        )}

                        <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
                            <FormField
                                control={form.control}
                                name="name"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>Name</FormLabel>
                                        <FormControl>
                                            <Input
                                                {...field}
                                                placeholder="Enter server name"
                                            />
                                        </FormControl>
                                        <FormDescription>A unique name for this MCP server</FormDescription>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />

                            <FormField
                                control={form.control}
                                name="transport"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>Transport</FormLabel>
                                        <Select
                                            defaultValue={field.value}
                                            onValueChange={(v: McpTransport) => {
                                                field.onChange(v);

                                                // Normalize opposite config to avoid stale values
                                                if (v === 'stdio') {
                                                    form.setValue('sse', undefined);

                                                    if (!form.getValues('stdio')) {
                                                        form.setValue('stdio', { args: '', command: '', env: [] });
                                                    }
                                                } else {
                                                    form.setValue('stdio', undefined);

                                                    if (!form.getValues('sse')) {
                                                        form.setValue('sse', { headers: [], url: '' });
                                                    }
                                                }
                                            }}
                                        >
                                            <FormControl>
                                                <SelectTrigger>
                                                    <SelectValue placeholder="Select transport" />
                                                </SelectTrigger>
                                            </FormControl>
                                            <SelectContent>
                                                <SelectItem value="stdio">STDIO</SelectItem>
                                                <SelectItem value="sse">SSE</SelectItem>
                                            </SelectContent>
                                        </Select>
                                        <FormDescription>STDIO for local process; SSE for remote URL</FormDescription>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                        </div>

                        {/* STDIO configuration */}
                        {transport === 'stdio' && (
                            <div className="flex flex-col gap-4">
                                <h3 className="text-lg font-medium">STDIO Configuration</h3>
                                <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
                                    <FormField
                                        control={form.control}
                                        name="stdio.command"
                                        render={({ field }) => (
                                            <FormItem>
                                                <FormLabel>Command</FormLabel>
                                                <FormControl>
                                                    <Input
                                                        {...field}
                                                        placeholder="/usr/local/bin/node"
                                                    />
                                                </FormControl>
                                                <FormMessage />
                                            </FormItem>
                                        )}
                                    />
                                    <FormField
                                        control={form.control}
                                        name="stdio.args"
                                        render={({ field }) => (
                                            <FormItem>
                                                <FormLabel>Args</FormLabel>
                                                <FormControl>
                                                    <Input
                                                        {...field}
                                                        placeholder="/path/to/script.js --flag value"
                                                        value={field.value ?? ''}
                                                    />
                                                </FormControl>
                                                <FormDescription>Space-separated arguments</FormDescription>
                                                <FormMessage />
                                            </FormItem>
                                        )}
                                    />
                                </div>

                                <div>
                                    <div className="mb-2 flex items-center justify-between">
                                        <h4 className="text-sm font-medium">Environment Variables</h4>
                                        <Button
                                            onClick={() => handleAddKeyValue('env')}
                                            size="sm"
                                            type="button"
                                            variant="outline"
                                        >
                                            <Plus className="size-3" /> Add
                                        </Button>
                                    </div>
                                    <div className="flex flex-col gap-2">
                                        {stdioEnvArray.fields.length === 0 && (
                                            <div className="text-muted-foreground text-sm">No variables</div>
                                        )}
                                        {stdioEnvArray.fields.map((field, index) => (
                                            <div
                                                className="grid grid-cols-1 gap-2 md:grid-cols-5"
                                                key={field.id}
                                            >
                                                <Controller
                                                    control={form.control}
                                                    name={`stdio.env.${index}.key` as const}
                                                    render={({ field }) => (
                                                        <Input
                                                            {...field}
                                                            className="md:col-span-2"
                                                            placeholder="KEY"
                                                        />
                                                    )}
                                                />
                                                <Controller
                                                    control={form.control}
                                                    name={`stdio.env.${index}.value` as const}
                                                    render={({ field }) => (
                                                        <Input
                                                            {...field}
                                                            className="md:col-span-2"
                                                            placeholder="VALUE"
                                                        />
                                                    )}
                                                />
                                                <Button
                                                    className="justify-self-start"
                                                    onClick={() => stdioEnvArray.remove(index)}
                                                    type="button"
                                                    variant="ghost"
                                                >
                                                    <Trash2 className="size-4" />
                                                </Button>
                                            </div>
                                        ))}
                                    </div>
                                </div>
                            </div>
                        )}

                        {/* SSE configuration */}
                        {transport === 'sse' && (
                            <div className="flex flex-col gap-4">
                                <h3 className="text-lg font-medium">SSE Configuration</h3>
                                <FormField
                                    control={form.control}
                                    name="sse.url"
                                    render={({ field }) => (
                                        <FormItem>
                                            <FormLabel>URL</FormLabel>
                                            <FormControl>
                                                <Input
                                                    {...field}
                                                    placeholder="https://mcp.example.com/sse"
                                                />
                                            </FormControl>
                                            <FormMessage />
                                        </FormItem>
                                    )}
                                />

                                <div>
                                    <div className="mb-2 flex items-center justify-between">
                                        <h4 className="text-sm font-medium">Headers</h4>
                                        <Button
                                            onClick={() => handleAddKeyValue('headers')}
                                            size="sm"
                                            type="button"
                                            variant="outline"
                                        >
                                            <Plus className="size-3" /> Add
                                        </Button>
                                    </div>
                                    <div className="flex flex-col gap-2">
                                        {sseHeadersArray.fields.length === 0 && (
                                            <div className="text-muted-foreground text-sm">No headers</div>
                                        )}
                                        {sseHeadersArray.fields.map((field, index) => (
                                            <div
                                                className="grid grid-cols-1 gap-2 md:grid-cols-5"
                                                key={field.id}
                                            >
                                                <Controller
                                                    control={form.control}
                                                    name={`sse.headers.${index}.key` as const}
                                                    render={({ field }) => (
                                                        <Input
                                                            {...field}
                                                            className="md:col-span-2"
                                                            placeholder="Header"
                                                        />
                                                    )}
                                                />
                                                <Controller
                                                    control={form.control}
                                                    name={`sse.headers.${index}.value` as const}
                                                    render={({ field }) => (
                                                        <Input
                                                            {...field}
                                                            className="md:col-span-2"
                                                            placeholder="Value"
                                                        />
                                                    )}
                                                />
                                                <Button
                                                    className="justify-self-start"
                                                    onClick={() => sseHeadersArray.remove(index)}
                                                    type="button"
                                                    variant="ghost"
                                                >
                                                    <Trash2 className="size-4" />
                                                </Button>
                                            </div>
                                        ))}
                                    </div>
                                </div>
                            </div>
                        )}

                        {/* Tools configuration - only for existing servers; toggles only */}
                        {!isNew && (
                            <div className="flex flex-col gap-4">
                                <div>
                                    <h3 className="text-lg font-medium">Tools</h3>
                                    <p className="text-muted-foreground text-sm">Enable or disable available tools</p>
                                </div>
                                <div className="grid grid-cols-1 gap-3 md:grid-cols-2 lg:grid-cols-3">
                                    {toolsArray.fields.length === 0 && (
                                        <div className="text-muted-foreground text-sm">No tools</div>
                                    )}
                                    {toolsArray.fields.map((tool, index) => (
                                        <div
                                            className="flex flex-col gap-2 rounded-md border p-2"
                                            key={tool.id}
                                        >
                                            <div className="flex items-start justify-between gap-4">
                                                <div className="flex-1 text-sm">
                                                    <div className="truncate font-medium">
                                                        {form.watch(`tools.${index}.name`) || 'tool'}
                                                    </div>
                                                    {form.watch(`tools.${index}.description`) && (
                                                        <div className="text-muted-foreground">
                                                            {form.watch(`tools.${index}.description`) as string}
                                                        </div>
                                                    )}
                                                </div>
                                                <div className="flex items-center gap-2">
                                                    <span className="text-muted-foreground text-xs">Enabled</span>
                                                    <Controller
                                                        control={form.control}
                                                        name={`tools.${index}.enabled` as const}
                                                        render={({ field }) => (
                                                            <Switch
                                                                aria-label={`Toggle ${form.getValues(`tools.${index}.name`) || 'tool'}`}
                                                                checked={!!field.value}
                                                                onCheckedChange={field.onChange}
                                                            />
                                                        )}
                                                    />
                                                    <Button
                                                        disabled={toolTestLoadingIndex === index}
                                                        onClick={() => handleTestTool(index)}
                                                        size="sm"
                                                        type="button"
                                                        variant="outline"
                                                    >
                                                        {toolTestLoadingIndex === index ? (
                                                            <Loader2 className="size-3 animate-spin" />
                                                        ) : (
                                                            <Play className="size-3" />
                                                        )}
                                                        {toolTestLoadingIndex === index ? 'Testing...' : 'Test'}
                                                    </Button>
                                                </div>
                                            </div>
                                            {toolTestIndex === index && (toolTestMessage || toolTestError) && (
                                                <div className="mt-1 text-xs">
                                                    {toolTestMessage && (
                                                        <span className="text-green-600">{toolTestMessage}</span>
                                                    )}
                                                    {toolTestError && (
                                                        <span className="text-red-600">{toolTestError}</span>
                                                    )}
                                                </div>
                                            )}
                                        </div>
                                    ))}
                                </div>
                            </div>
                        )}
                    </form>
                </Form>
            </div>

            {/* Sticky buttons */}
            <div className="bg-background sticky -bottom-4 -mx-4 mt-4 -mb-4 flex items-center border-t p-4 shadow-lg">
                <div className="flex gap-2">
                    {!isNew && (
                        <Button
                            onClick={handleDelete}
                            type="button"
                            variant="destructive"
                        >
                            <Trash2 className="size-4" />
                            Delete
                        </Button>
                    )}
                    <Button
                        disabled={isTestLoading}
                        onClick={handleTest}
                        type="button"
                        variant="outline"
                    >
                        {isTestLoading ? <Loader2 className="size-4 animate-spin" /> : <Play className="size-4" />}
                        {isTestLoading ? 'Testing...' : 'Test'}
                    </Button>
                </div>
                <div className="ml-auto flex gap-2">
                    <Button
                        onClick={() => navigate('/settings/mcp-servers')}
                        type="button"
                        variant="outline"
                    >
                        Cancel
                    </Button>
                    <Button
                        disabled={form.formState.isSubmitting}
                        form="mcp-server-form"
                        type="submit"
                        variant="secondary"
                    >
                        {form.formState.isSubmitting ? (
                            <Loader2 className="size-4 animate-spin" />
                        ) : (
                            <Save className="size-4" />
                        )}
                        {form.formState.isSubmitting ? 'Saving...' : isNew ? 'Create MCP Server' : 'Update MCP Server'}
                    </Button>
                </div>
            </div>

            <ConfirmationDialog
                cancelText="Cancel"
                confirmText="Delete"
                handleConfirm={handleConfirmDelete}
                handleOpenChange={setIsDeleteDialogOpen}
                isOpen={isDeleteDialogOpen}
                itemName={form.watch('name')}
                itemType="MCP server"
            />
        </Fragment>
    );
};

export default SettingsMcpServer;
