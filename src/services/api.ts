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
        const response = await api.get<Trip[]>('/partner/rfqs');
        return response.data;
    },

    listRFQHistory: async (): Promise<Trip[]> => {
        const response = await api.get<Trip[]>('/partner/rfqs/history');
        return response.data;
    },

    listAssignedTrips: async (): Promise<Trip[]> => {
        const response = await api.get<Trip[]>('/partner/trips');
        return response.data;
    },

    submitQuote: async (req: QuoteRequest) => {
        const response = await api.post('/partner/quotes', req);
        return response.data;
    },

    assignDriver: async (req: AssignDriverRequest) => {
        const response = await api.post('/partner/assign', req);
        return response.data;
    },

    acceptRFQ: async (tripId: string) => {
        const response = await api.post('/partner/accept', { trip_id: tripId });
        return response.data;
    },

    rejectRFQ: async (tripId: string) => {
        const response = await api.post('/partner/reject', { trip_id: tripId });
        return response.data;
    },

    listUninvoicedTrips: async (): Promise<Trip[]> => {
        const response = await api.get<Trip[]>('/partner/trips/completed');
        return response.data;
    },

    listInvoices: async (): Promise<Invoice[]> => {
        const response = await api.get<Invoice[]>('/partner/invoices');
        return response.data;
    },

    submitInvoice: async (data: { trip_id: string; invoice_number: string; amount: number }) => {
        const response = await api.post('/partner/invoices', data);
        return response.data;
    },

    updateInvoice: async (invoiceId: string, data: { invoice_number: string; amount: number }) => {
        const response = await api.put(`/partner/invoices/${invoiceId}`, data);
        return response.data;
    },

    closeInvoice: async (invoiceId: string) => {
        const response = await api.post(`/partner/invoices/${invoiceId}/close`);
        return response.data;
    },

    listVehicles: async (): Promise<Vehicle[]> => {
        const response = await api.get<Vehicle[]>('/partner/vehicles');
        return response.data;
    },

    createVehicle: async (data: Partial<Vehicle>) => {
        const response = await api.post('/partner/vehicles', data);
        return response.data;
    },

    listDrivers: async (): Promise<Driver[]> => {
        const response = await api.get<Driver[]>('/partner/drivers');
        return response.data;
    },

    createDriver: async (data: Partial<Driver>) => {
        const response = await api.post('/partner/drivers', data);
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
}

export interface Driver {
    id: string;
    name: string;
    phone: string;
    license_number: string;
    current_vehicle_id?: string;
    license_expiry?: string;
    itc_permit_expiry?: string;
    visa_expiry?: string;
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
}

export const operatorService = {
    listOutsourceCompanies: async () => {
        const response = await api.get('/operator/outsource-companies');
        return response.data;
    },

    createOutsourceCompany: async (data: any) => {
        const response = await api.post('/operator/outsource-companies', data);
        return response.data;
    },

    updateOutsourceCompany: async (id: string, data: any) => {
        const response = await api.put(`/operator/outsource-companies/${id}`, data);
        return response.data;
    },

    deleteOutsourceCompany: async (id: string) => {
        const response = await api.delete(`/operator/outsource-companies/${id}`);
        return response.data;
    },

    getSubmittedQuotes: async (): Promise<TripOffer[]> => {
        const response = await api.get<TripOffer[]>('/operator/quotes');
        return response.data;
    },

    acceptQuote: async (quoteId: string) => {
        const response = await api.post(`/operator/quotes/${quoteId}/accept`, {});
        return response.data;
    },

    rejectQuote: async (quoteId: string) => {
        const response = await api.post(`/operator/quotes/${quoteId}/reject`, {});
        return response.data;
    },

    listOperatorTrips: async (): Promise<Trip[]> => {
        const response = await api.get<Trip[]>('/operator/trips');
        return response.data;
    },

    assignOutsource: async (tripId: string, partnerIds?: string[]) => {
        const response = await api.post(`/operator/trips/${tripId}/assign`, {
            partner_ids: partnerIds
        });
        return response.data;
    },

    createTrip: async (data: Partial<Trip>) => {
        const response = await api.post('/operator/trips', data);
        return response.data;
    },

    listVehicles: async () => {
        const response = await api.get('/operator/vehicles');
        return response.data;
    },

    listDrivers: async () => {
        const response = await api.get('/operator/drivers');
        return response.data;
    },

    dispatchTrip: async (tripId: string, driverId: string, vehicleId: string) => {
        const response = await api.post(`/operator/trips/${tripId}/dispatch`, {
            driver_id: driverId,
            vehicle_id: vehicleId
        });
        return response.data;
    }
};

export const tripService = {
    getTrip: async (id: string): Promise<Trip> => {
        const response = await api.get<Trip>(`/trips/${id}`);
        return response.data;
    }
};

export const authService = {
    changePassword: async (userId: string, newPassword: string) => {
        const response = await api.post('/auth/change-password', { user_id: userId, new_password: newPassword });
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

