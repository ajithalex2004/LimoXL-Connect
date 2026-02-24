import { useState, useEffect } from 'react';
import { partnerService, type Trip, type Invoice } from '../services/api';
import { Check, AlertCircle, Plus } from 'lucide-react';

const Invoices = () => {
    const [activeTab, setActiveTab] = useState<'pending' | 'submitted'>('pending');
    const [trips, setTrips] = useState<Trip[]>([]);
    const [invoices, setInvoices] = useState<Invoice[]>([]);
    const [loading, setLoading] = useState(true);
    const [showModal, setShowModal] = useState(false);
    const [selectedTrip, setSelectedTrip] = useState<Trip | null>(null);
    const [selectedInvoice, setSelectedInvoice] = useState<Invoice | null>(null);
    const [isEditMode, setIsEditMode] = useState(false);
    const [invoiceForm, setInvoiceForm] = useState({
        invoiceNumber: '',
        amount: 0
    });

    useEffect(() => {
        loadData();
    }, [activeTab]);

    const loadData = async () => {
        setLoading(true);
        try {
            if (activeTab === 'pending') {
                const data = await partnerService.listUninvoicedTrips();
                setTrips(data);
            } else {
                const data = await partnerService.listInvoices();
                setInvoices(data);
            }
        } catch (error) {
            console.error("Failed to load invoice data", error);
        } finally {
            setLoading(false);
        }
    };

    const handleOpenSubmit = (trip: Trip) => {
        setSelectedTrip(trip);
        setSelectedInvoice(null);
        setIsEditMode(false);
        setInvoiceForm({
            invoiceNumber: '',
            amount: trip.price || 0
        });
        setShowModal(true);
    };

    const handleOpenEdit = (invoice: Invoice) => {
        setSelectedInvoice(invoice);
        setSelectedTrip(null);
        setIsEditMode(true);
        setInvoiceForm({
            invoiceNumber: invoice.invoice_number,
            amount: invoice.amount
        });
        setShowModal(true);
    };

    const handleCloseInvoice = async (invoiceId: string) => {
        if (!confirm('Are you sure you want to close this invoice? This action cannot be undone.')) {
            return;
        }

        try {
            await partnerService.closeInvoice(invoiceId);
            loadData();
        } catch (error) {
            console.error("Failed to close invoice", error);
            alert("Failed to close invoice.");
        }
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        try {
            if (isEditMode && selectedInvoice) {
                // Update existing invoice
                await partnerService.updateInvoice(selectedInvoice.id, {
                    invoice_number: invoiceForm.invoiceNumber,
                    amount: invoiceForm.amount
                });
            } else if (selectedTrip) {
                // Submit new invoice
                await partnerService.submitInvoice({
                    trip_id: selectedTrip.id,
                    invoice_number: invoiceForm.invoiceNumber,
                    amount: invoiceForm.amount
                });
            }
            setShowModal(false);
            setActiveTab('submitted');
            loadData();
        } catch (error) {
            console.error("Failed to submit/update invoice", error);
            alert("Failed to submit/update invoice.");
        }
    };

    return (
        <div className="space-y-6">
            <div>
                <h1 className="text-2xl font-bold text-gray-900">Invoices</h1>
                <p className="text-gray-500">Manage trip invoicing and payments</p>
            </div>

            {/* Tabs */}
            <div className="flex border-b border-gray-200">
                <button
                    onClick={() => setActiveTab('pending')}
                    className={`pb-4 px-4 text-sm font-medium transition-colors relative ${activeTab === 'pending'
                        ? 'text-blue-600 border-b-2 border-blue-600'
                        : 'text-gray-500 hover:text-gray-700'
                        }`}
                >
                    Pending Submission
                </button>
                <button
                    onClick={() => setActiveTab('submitted')}
                    className={`pb-4 px-4 text-sm font-medium transition-colors relative ${activeTab === 'submitted'
                        ? 'text-blue-600 border-b-2 border-blue-600'
                        : 'text-gray-500 hover:text-gray-700'
                        }`}
                >
                    Submitted Invoices
                </button>
            </div>

            {/* Content */}
            {loading ? (
                <div className="text-center py-12 text-gray-500">Loading...</div>
            ) : (
                <>
                    {activeTab === 'pending' ? (
                        <div className="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
                            {trips.length === 0 ? (
                                <div className="p-8 text-center text-gray-500">No completed trips pending invoice.</div>
                            ) : (
                                <table className="w-full text-left text-sm">
                                    <thead className="bg-gray-50 border-b border-gray-200">
                                        <tr>
                                            <th className="px-6 py-3 font-semibold text-gray-900">Trip Ref</th>
                                            <th className="px-6 py-3 font-semibold text-gray-900">Date/Time</th>
                                            <th className="px-6 py-3 font-semibold text-gray-900">Route</th>
                                            <th className="px-6 py-3 font-semibold text-gray-900">Est. Price</th>
                                            <th className="px-6 py-3 font-semibold text-gray-900 text-right">Action</th>
                                        </tr>
                                    </thead>
                                    <tbody className="divide-y divide-gray-100">
                                        {trips.map((trip) => (
                                            <tr key={trip.id} className="hover:bg-gray-50">
                                                <td className="px-6 py-4 font-medium text-gray-900">{trip.booking_reference}</td>
                                                <td className="px-6 py-4 text-gray-600">
                                                    {new Date(trip.pickup_time).toLocaleDateString()} {new Date(trip.pickup_time).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                                                </td>
                                                <td className="px-6 py-4 text-gray-600">
                                                    <span className="font-medium">{trip.pickup_zone}</span> → {trip.dropoff_zone}
                                                </td>
                                                <td className="px-6 py-4 text-gray-600">
                                                    AED {trip.price?.toFixed(2) || '0.00'}
                                                </td>
                                                <td className="px-6 py-4 text-right">
                                                    <button
                                                        onClick={() => handleOpenSubmit(trip)}
                                                        className="text-blue-600 hover:text-blue-800 font-medium flex items-center justify-end gap-1 w-full"
                                                    >
                                                        <Plus className="h-4 w-4" /> Submit Invoice
                                                    </button>
                                                </td>
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>
                            )}
                        </div>
                    ) : (
                        <div className="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
                            {invoices.length === 0 ? (
                                <div className="p-8 text-center text-gray-500">No invoices submitted yet.</div>
                            ) : (
                                <table className="w-full text-left text-sm">
                                    <thead className="bg-gray-50 border-b border-gray-200">
                                        <tr>
                                            <th className="px-6 py-3 font-semibold text-gray-900">Invoice #</th>
                                            <th className="px-6 py-3 font-semibold text-gray-900">Trip Ref</th>
                                            <th className="px-6 py-3 font-semibold text-gray-900">Gross</th>
                                            <th className="px-6 py-3 font-semibold text-gray-900">Fee (10%)</th>
                                            <th className="px-6 py-3 font-semibold text-gray-900">VAT (5%)</th>
                                            <th className="px-6 py-3 font-semibold text-gray-900">Net Payout</th>
                                            <th className="px-6 py-3 font-semibold text-gray-900">Status</th>
                                            <th className="px-6 py-3 font-semibold text-gray-900 text-right">Actions</th>
                                        </tr>
                                    </thead>
                                    <tbody className="divide-y divide-gray-100">
                                        {invoices.map((inv) => (
                                            <tr key={inv.id} className="hover:bg-gray-50">
                                                <td className="px-6 py-4 font-medium text-gray-900">{inv.invoice_number}</td>
                                                <td className="px-6 py-4 text-gray-600">{inv.booking_reference}</td>
                                                <td className="px-6 py-4 text-gray-600 font-medium">AED {inv.amount.toFixed(2)}</td>
                                                <td className="px-6 py-4 text-red-500 font-medium whitespace-nowrap">- AED {inv.platform_fee?.toFixed(2) || '0.00'}</td>
                                                <td className="px-6 py-4 text-blue-500 font-medium whitespace-nowrap">+ AED {inv.vat_amount?.toFixed(2) || '0.00'}</td>
                                                <td className="px-6 py-4 font-bold text-gray-900 bg-emerald-50/50 whitespace-nowrap">
                                                    AED {inv.net_payout?.toFixed(2) || inv.amount.toFixed(2)}
                                                </td>
                                                <td className="px-6 py-4">
                                                    <span className={`inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium border ${inv.status === 'PAID'
                                                        ? 'bg-green-50 text-green-700 border-green-200'
                                                        : inv.status === 'REJECTED'
                                                            ? 'bg-red-50 text-red-700 border-red-200'
                                                            : inv.status === 'CLOSED'
                                                                ? 'bg-gray-50 text-gray-700 border-gray-200'
                                                                : 'bg-yellow-50 text-yellow-700 border-yellow-200'
                                                        }`}>
                                                        {inv.status === 'PAID' && <Check className="h-3 w-3" />}
                                                        {inv.status === 'REJECTED' && <AlertCircle className="h-3 w-3" />}
                                                        {inv.status}
                                                    </span>
                                                </td>
                                                <td className="px-6 py-4 text-right">
                                                    <div className="flex gap-2 justify-end">
                                                        {inv.status === 'PENDING' && (
                                                            <>
                                                                <button
                                                                    onClick={() => handleOpenEdit(inv)}
                                                                    className="text-blue-600 hover:text-blue-800 text-sm font-medium"
                                                                >
                                                                    Resubmit
                                                                </button>
                                                                <button
                                                                    onClick={() => handleCloseInvoice(inv.id)}
                                                                    className="text-gray-600 hover:text-gray-800 text-sm font-medium"
                                                                >
                                                                    Close
                                                                </button>
                                                            </>
                                                        )}
                                                        {inv.status === 'CLOSED' && (
                                                            <span className="text-gray-400 text-sm">No actions</span>
                                                        )}
                                                    </div>
                                                </td>
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>
                            )}
                        </div>
                    )}
                </>
            )}

            {/* Modal */}
            {showModal && (selectedTrip || selectedInvoice) && (
                <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm">
                    <div className="bg-white rounded-xl shadow-xl max-w-md w-full p-6">
                        <h3 className="text-lg font-bold text-gray-900 mb-2">
                            {isEditMode ? 'Resubmit Invoice' : 'Submit Invoice'}
                        </h3>
                        <p className="text-sm text-gray-500 mb-4">
                            {isEditMode ? `Invoice: ${selectedInvoice?.invoice_number}` : `Trip: ${selectedTrip?.booking_reference}`}
                        </p>

                        <form onSubmit={handleSubmit} className="space-y-4">
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">Invoice Number</label>
                                <input
                                    type="text"
                                    required
                                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                                    placeholder="e.g. INV-2024-001"
                                    value={invoiceForm.invoiceNumber}
                                    onChange={e => setInvoiceForm({ ...invoiceForm, invoiceNumber: e.target.value })}
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">Amount (AED)</label>
                                <input
                                    type="number"
                                    step="0.01"
                                    required
                                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                                    value={invoiceForm.amount}
                                    onChange={e => setInvoiceForm({ ...invoiceForm, amount: parseFloat(e.target.value) })}
                                />
                            </div>

                            {invoiceForm.amount > 0 && (
                                <div className="bg-gray-50 p-4 rounded-lg space-y-2 border border-gray-200">
                                    <div className="flex justify-between text-sm text-gray-600">
                                        <span>Gross Amount:</span>
                                        <span>AED {invoiceForm.amount.toFixed(2)}</span>
                                    </div>
                                    <div className="flex justify-between text-sm text-red-600">
                                        <span>Platform Fee (10%):</span>
                                        <span>- AED {(invoiceForm.amount * 0.1).toFixed(2)}</span>
                                    </div>
                                    <div className="flex justify-between text-sm text-emerald-600">
                                        <span>VAT Amount (5%):</span>
                                        <span>+ AED {(invoiceForm.amount * 0.05).toFixed(2)}</span>
                                    </div>
                                    <div className="flex justify-between font-bold text-gray-900 border-t pt-2 mt-2">
                                        <span>Estimated Payout:</span>
                                        <span>AED {(invoiceForm.amount - (invoiceForm.amount * 0.1) + (invoiceForm.amount * 0.05)).toFixed(2)}</span>
                                    </div>
                                </div>
                            )}
                            <div className="flex gap-3 pt-4">
                                <button
                                    type="button"
                                    onClick={() => setShowModal(false)}
                                    className="flex-1 px-4 py-2 text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200"
                                >
                                    Cancel
                                </button>
                                <button
                                    type="submit"
                                    className="flex-1 px-4 py-2 text-white bg-blue-600 rounded-lg hover:bg-blue-700"
                                >
                                    Submit Invoice
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            )}
        </div>
    );
};

export default Invoices;
