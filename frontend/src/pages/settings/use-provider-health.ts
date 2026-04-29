import { useState } from 'react';

import type { ProviderType } from '@/graphql/types';

import type { ProviderHealthResponse } from './provider-form-schema';

export interface EndpointCheckPayload {
    apiKey?: string;
    apiUrl?: string;
    embeddingModel?: string;
    type: ProviderType;
}

export function useProviderHealth() {
    const [isEndpointHealthLoading, setIsEndpointHealthLoading] = useState(false);
    const [endpointHealth, setEndpointHealth] = useState<null | ProviderHealthResponse>(null);

    const checkEndpoint = async (payload: EndpointCheckPayload): Promise<void> => {
        setIsEndpointHealthLoading(true);

        try {
            const resp = await fetch('/api/v1/providers/health', {
                body: JSON.stringify(payload),
                credentials: 'include',
                headers: {
                    'Content-Type': 'application/json',
                },
                method: 'POST',
            });

            const body = await resp.json();
            if (!resp.ok || body?.status !== 'success') {
                throw new Error(body?.msg || body?.error || 'Failed to check provider endpoint health');
            }

            setEndpointHealth(body?.data as ProviderHealthResponse);
        } catch (error) {
            setEndpointHealth(null);
            throw error;
        } finally {
            setIsEndpointHealthLoading(false);
        }
    };

    return { checkEndpoint, endpointHealth, isEndpointHealthLoading };
}
