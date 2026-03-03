import { api } from '../lib/auth';

export interface Trip {
    id: string;
    booking_reference: string;
    pickup_zone: string;
    dropoff_zone: string;
    pickup_time: string;
    status: string;
    passenger_name: string;
    passenger_phone: string;
    price: number;
    payment_method: string;
    supplier_name: string;
    service_type?: string;
    requested_vehicle_type?: string;
    requested_vehicle_class?: string;
    requested_vehicle_group?: string;
    rfq_number?: string;
    pickup_landmark?: string;
    dropoff_landmark?: string;
    driver_link_token?: string;
    created_at?: string;
    updated_at?: string;
}

export interface BookingConfig {
    id: string;
    company_id: string;
    name: string;
    booking_type: string;
    request_type: string;
    priority: string;
    sort_order: number;
    vehicle_classes: string[];
    vehicle_groups: string[];
    vehicle_usages: string[];
    pickup_buffer: number;
    auto_dispatch_buffer: number;
    pricing_source: string;
    approval_workflow_required: boolean;
    epod_required: boolean;
    is_active: boolean;
    created_at?: string;
    updated_at?: string;
}

export interface QuoteRequest {
    trip_id: string;
    price: number;
    notes?: string;
}

export interface AssignDriverRequest {
    trip_id: string;
    driver_name: string;
    driver_phone: string;
    vehicle_model: string;
    vehicle_plate: string;
}

export const partnerService = {
    listRFQs: async (): Promise<Trip[]> => {
        const response = await api.get<Trip[]>('partner/rfqs');
        return response.data;
    },

    listRFQHistory: async (): Promise<Trip[]> => {
        const response = await api.get<Trip[]>('partner/rfqs/history');
        return response.data;
    },

    listAssignedTrips: async (): Promise<Trip[]> => {
        const response = await api.get<Trip[]>('partner/trips');
        return response.data;
    },

    submitQuote: async (req: QuoteRequest) => {
        const response = await api.post('partner/quotes', req);
        return response.data;
    },

    assignDriver: async (req: AssignDriverRequest) => {
        const response = await api.post('partner/assign', req);
        return response.data;
    },

    acceptRFQ: async (tripId: string) => {
        const response = await api.post('partner/accept', { trip_id: tripId });
        return response.data;
    },

    rejectRFQ: async (tripId: string) => {
        const response = await api.post('partner/reject', { trip_id: tripId });
        return response.data;
    },

    listUninvoicedTrips: async (): Promise<Trip[]> => {
        const response = await api.get<Trip[]>('partner/trips/completed');
        return response.data;
    },

    listInvoices: async (): Promise<Invoice[]> => {
        const response = await api.get<Invoice[]>('partner/invoices');
        return response.data;
    },

    submitInvoice: async (data: { trip_id: string; invoice_number: string; amount: number }) => {
        const response = await api.post('partner/invoices', data);
        return response.data;
    },

    updateInvoice: async (invoiceId: string, data: { invoice_number: string; amount: number }) => {
        const response = await api.put(`partner/invoices/${invoiceId}`, data);
        return response.data;
    },

    closeInvoice: async (invoiceId: string) => {
        const response = await api.post(`partner/invoices/${invoiceId}/close`);
        return response.data;
    },

    listVehicles: async (): Promise<Vehicle[]> => {
        const response = await api.get<Vehicle[]>('partner/vehicles');
        return response.data;
    },

    createVehicle: async (data: Partial<Vehicle>) => {
        const response = await api.post('partner/vehicles', data);
        return response.data;
    },

    listDrivers: async (): Promise<Driver[]> => {
        const response = await api.get<Driver[]>('partner/drivers');
        return response.data;
    },

    createDriver: async (data: Partial<Driver>) => {
        const response = await api.post('partner/drivers', data);
        return response.data;
    }
};

export interface Vehicle {
    id: string;
    company_id: string;
    plate_number: string;
    type: string;
    vehicle_class: string;
    vehicle_group: string;
    model: string;
    capacity: number;
    status: string;
    permit_expiry?: string;
    insurance_expiry?: string;
    chassis_no?: string;
    vin?: string;
    year_of_manufacture?: number;
    color?: string;
    registration_number?: string;
    plate_code?: string;
    plate_category?: string;
    emirate?: string;
    hierarchy?: string;
    vehicle_usage?: string;
    created_at?: string;
    updated_at?: string;
}

export interface Driver {
    id: string;
    company_id?: string;
    name: string;
    phone: string;
    license_number: string;
    current_vehicle_id?: string;
    license_expiry?: string;
    itc_permit_expiry?: string;
    visa_expiry?: string;
    first_name?: string;
    last_name?: string;
    date_of_birth?: string;
    nationality?: string;
    emirates_id?: string;
    date_of_join?: string;
    dallas_id?: string;
    communication_language?: string;
    hierarchy?: string;
    driver_type?: string;
    created_at?: string;
}

export interface FleetAttachment {
    id: string;
    entity_id: string;
    entity_type: string;
    file_name: string;
    file_url: string;
    file_type: string;
    created_at?: string;
}

export interface NUIMaster {
    id: string;
    company_id: string;
    category: 'TYPE' | 'CLASS' | 'USAGE' | 'DRIVER_TYPE' | 'HIERARCHY';
    name: string;
    description: string;
    is_active: boolean;
    created_at?: string;
    updated_at?: string;
}

export interface TripOffer {
    id: string;
    trip_id: string;
    supplier_company_id: string;
    supplier_name: string;
    status: 'PENDING' | 'ACCEPTED' | 'REJECTED' | 'TIMEOUT';
    score: number;
    price: number;
    notes?: string;
    created_at: string;
    booking_reference?: string;
    rfq_number?: string;
    service_type?: string;
}

export const operatorService = {
    listOutsourceCompanies: async () => {
        const response = await api.get('operator/outsource-companies');
        return response.data;
    },

    createOutsourceCompany: async (data: any) => {
        const response = await api.post('operator/outsource-companies', data);
        return response.data;
    },

    updateOutsourceCompany: async (id: string, data: any) => {
        const response = await api.put(`operator/outsource-companies/${id}`, data);
        return response.data;
    },

    deleteOutsourceCompany: async (id: string) => {
        const response = await api.delete(`operator/outsource-companies/${id}`);
        return response.data;
    },

    getSubmittedQuotes: async (): Promise<TripOffer[]> => {
        const response = await api.get<TripOffer[]>('operator/quotes');
        return response.data;
    },

    acceptQuote: async (quoteId: string) => {
        const response = await api.post(`operator/quotes/${quoteId}/accept`, {});
        return response.data;
    },

    rejectQuote: async (quoteId: string) => {
        const response = await api.post(`operator/quotes/${quoteId}/reject`, {});
        return response.data;
    },

    listOperatorTrips: async (): Promise<Trip[]> => {
        const response = await api.get<Trip[]>('operator/trips');
        return response.data;
    },

    assignOutsource: async (tripId: string, partnerIds?: string[]) => {
        const response = await api.post(`operator/trips/${tripId}/assign`, {
            partner_ids: partnerIds
        });
        return response.data;
    },

    createTrip: async (data: Partial<Trip>) => {
        const response = await api.post('operator/trips', data);
        return response.data;
    },

    listVehicles: async (): Promise<Vehicle[]> => {
        const response = await api.get<Vehicle[]>('operator/vehicles');
        return response.data;
    },

    createVehicle: async (data: Partial<Vehicle>) => {
        const response = await api.post('operator/vehicles', data);
        return response.data;
    },

    updateVehicle: async (id: string, data: Partial<Vehicle>) => {
        const response = await api.put(`operator/vehicles/${id}`, data);
        return response.data;
    },

    deleteVehicle: async (id: string) => {
        const response = await api.delete(`operator/vehicles/${id}`);
        return response.data;
    },

    listDrivers: async (): Promise<Driver[]> => {
        const response = await api.get<Driver[]>('operator/drivers');
        return response.data;
    },

    createDriver: async (data: Partial<Driver>) => {
        const response = await api.post('operator/drivers', data);
        return response.data;
    },

    updateDriver: async (id: string, data: Partial<Driver>) => {
        const response = await api.put(`operator/drivers/${id}`, data);
        return response.data;
    },

    deleteDriver: async (id: string) => {
        const response = await api.delete(`operator/drivers/${id}`);
        return response.data;
    },

    listAttachments: async (entityID: string, entityType: string): Promise<FleetAttachment[]> => {
        const response = await api.get<FleetAttachment[]>('operator/attachments', {
            params: { entity_id: entityID, entity_type: entityType }
        });
        return response.data;
    },

    createAttachment: async (data: Partial<FleetAttachment>) => {
        const response = await api.post('operator/attachments', data);
        return response.data;
    },

    deleteAttachment: async (id: string) => {
        const response = await api.delete(`operator/attachments/${id}`);
        return response.data;
    },

    dispatchTrip: async (tripId: string, driverId: string, vehicleId: string) => {
        const response = await api.post(`operator/trips/${tripId}/dispatch`, {
            driver_id: driverId,
            vehicle_id: vehicleId
        });
        return response.data;
    },

    listUsers: async () => {
        const response = await api.get('operator/users');
        return response.data;
    },

    createTeamMember: async (data: any) => {
        const response = await api.post('operator/users', data);
        return response.data;
    },

    // NUI Masters
    listMasters: async (category?: string): Promise<NUIMaster[]> => {
        const response = await api.get<NUIMaster[]>('operator/masters', { params: { category } });
        return response.data || [];
    },

    createMaster: async (data: Partial<NUIMaster>) => {
        const response = await api.post('operator/masters', data);
        return response.data;
    },

    updateMaster: async (id: string, data: Partial<NUIMaster>) => {
        const response = await api.put(`operator/masters/${id}`, data);
        return response.data;
    },

    deleteMaster: async (id: string) => {
        const response = await api.delete(`operator/masters/${id}`);
        return response.data;
    },

    // Booking Configs
    listBookingConfigs: async (): Promise<BookingConfig[]> => {
        const response = await api.get<BookingConfig[]>('operator/booking-configs');
        return response.data || [];
    },

    createBookingConfig: async (data: Partial<BookingConfig>) => {
        const response = await api.post('operator/booking-configs', data);
        return response.data;
    },

    updateBookingConfig: async (id: string, data: Partial<BookingConfig>) => {
        const response = await api.put(`operator/booking-configs/${id}`, data);
        return response.data;
    },

    deleteBookingConfig: async (id: string) => {
        const response = await api.delete(`operator/booking-configs/${id}`);
        return response.data;
    }
};

export const tripService = {
    getTrip: async (id: string): Promise<Trip> => {
        const response = await api.get<Trip>(`trips/${id}`);
        return response.data;
    }
};

export const authService = {
    changePassword: async (userId: string, newPassword: string) => {
        const response = await api.post('auth/change-password', { user_id: userId, new_password: newPassword });
        return response.data;
    }
};

export interface Invoice {
    id: string;
    trip_id: string;
    supplier_company_id: string;
    invoice_number: string;
    amount: number;
    platform_fee?: number;
    vat_amount?: number;
    net_payout?: number;
    status: 'PENDING' | 'PAID' | 'REJECTED' | 'CLOSED';
    created_at: string;
    booking_reference?: string;
}

export interface SubmitInvoiceRequest {
    trip_id: string;
    invoice_number: string;
    amount: number;
}

