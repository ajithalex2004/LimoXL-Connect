import { api } from '../lib/auth';

export interface Tenant {
    id: string;
    company_id: string;
    name: string;
    slug: string;
    status: 'ACTIVE' | 'SUSPENDED' | 'TRIAL';
    plan: 'STARTER' | 'PROFESSIONAL' | 'ENTERPRISE';
    max_users: number;
    created_at: string;
}

export interface TenantFeature {
    id: string;
    tenant_id: string;
    feature_key: string;
    is_enabled: boolean;
    config: any;
}

export interface TenantWithFeatures extends Tenant {
    features: TenantFeature[];
    user_count: number;
    trip_count: number;
    company_name: string;
}

export const superAdminService = {
    listTenants: async (): Promise<TenantWithFeatures[]> => {
        const response = await api.get<TenantWithFeatures[]>('superadmin/tenants');
        return response.data;
    },

    createTenant: async (data: Partial<Tenant>) => {
        const response = await api.post('superadmin/tenants', data);
        return response.data;
    },

    updateTenant: async (id: string, data: Partial<Tenant>) => {
        const response = await api.put(`superadmin/tenants/${id}`, data);
        return response.data;
    },

    toggleFeature: async (tenantId: string, featureKey: string, isEnabled: boolean) => {
        const response = await api.post(`superadmin/tenants/${tenantId}/features`, {
            feature_key: featureKey,
            is_enabled: isEnabled
        });
        return response.data;
    }
};
