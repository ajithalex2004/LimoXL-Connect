import { useState, useEffect } from 'react';
import { operatorService, type NUIMaster } from '../../services/api';
import {
    Plus, Search, Edit2, Trash2, X,
    Settings2, Layers, Tag as TagIcon, Box, User
} from 'lucide-react';

const NUIMasters = () => {
    const [masters, setMasters] = useState<NUIMaster[]>([]);
    const [loading, setLoading] = useState(true);
    const [searchTerm, setSearchTerm] = useState('');
    const [activeCategory, setActiveCategory] = useState<'TYPE' | 'CLASS' | 'USAGE' | 'DRIVER_TYPE' | 'HIERARCHY'>('TYPE');

    // Modal State
    const [showModal, setShowModal] = useState(false);
    const [editingMaster, setEditingMaster] = useState<NUIMaster | null>(null);
    const [formData, setFormData] = useState<Partial<NUIMaster>>({
        name: '',
        description: '',
        is_active: true,
        category: 'TYPE'
    });

    useEffect(() => {
        loadData();
    }, [activeCategory]);

    const loadData = async () => {
        setLoading(true);
        try {
            const data = await operatorService.listMasters(activeCategory);
            setMasters(data || []);
        } catch (error) {
            console.error("Failed to load masters", error);
        } finally {
            setLoading(false);
        }
    };

    const handleOpenAdd = () => {
        setEditingMaster(null);
        setFormData({
            name: '',
            description: '',
            is_active: true,
            category: activeCategory
        });
        setShowModal(true);
    };

    const handleOpenEdit = (master: NUIMaster) => {
        setEditingMaster(master);
        setFormData({ ...master });
        setShowModal(true);
    };

    const handleDelete = async (id: string) => {
        if (!confirm('Are you sure you want to delete this master entry?')) return;
        try {
            await operatorService.deleteMaster(id);
            loadData();
        } catch (error) {
            console.error("Failed to delete master", error);
            alert("Failed to delete. It might be in use.");
        }
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        console.log("DEBUG: Submitting NUI Master:", { ...formData, category: activeCategory });
        try {
            if (editingMaster) {
                await operatorService.updateMaster(editingMaster.id, formData);
            } else {
                await operatorService.createMaster({ ...formData, category: activeCategory });
            }
            setShowModal(false);
            loadData();
        } catch (error) {
            console.error("Failed to save master", error);
            alert("Failed to save. Check for duplicates.");
        }
    };

    const filteredMasters = masters.filter(m =>
        m.name.toLowerCase().includes(searchTerm.toLowerCase())
    );

    const categories = [
        { id: 'TYPE', name: 'Vehicle Type', icon: Box, desc: 'Manage sedan, SUV, Van categories' },
        { id: 'CLASS', name: 'Vehicle Class', icon: Layers, desc: 'Define Luxury, Standard, Economy tiers' },
        { id: 'USAGE', name: 'Vehicle Usage', icon: TagIcon, desc: 'Categorize by Public or Private usage' },
        { id: 'DRIVER_TYPE', name: 'Driver Type', icon: User, desc: 'Regular, Temporary, VIP categories' },
        { id: 'HIERARCHY', name: 'Hierarchy', icon: Settings2, desc: 'Alpha, Beta, Executive levels' }
    ];

    return (
        <div className="space-y-6">
            <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
                <div>
                    <h1 className="text-2xl font-bold text-gray-900">NUI Masters</h1>
                    <p className="text-gray-500 text-sm">Configure dynamic dropdowns and categories for your fleet</p>
                </div>
                <div className="flex items-center gap-3">
                    <div className="relative">
                        <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
                        <input
                            type="text"
                            placeholder="Search masters..."
                            className="pl-10 pr-4 py-2 border border-gray-200 rounded-lg text-sm w-64 focus:outline-none focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500"
                            value={searchTerm}
                            onChange={e => setSearchTerm(e.target.value)}
                        />
                    </div>
                    <button
                        onClick={handleOpenAdd}
                        className="flex items-center gap-2 px-4 py-2 bg-emerald-600 text-white rounded-lg hover:bg-emerald-700 transition-colors text-sm font-semibold shadow-sm"
                    >
                        <Plus className="h-4 w-4" /> Add New
                    </button>
                </div>
            </div>

            {/* Category Selector */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                {categories.map((cat) => (
                    <button
                        key={cat.id}
                        onClick={() => setActiveCategory(cat.id as any)}
                        className={`p-4 rounded-xl border text-left transition-all ${activeCategory === cat.id
                            ? 'bg-emerald-50 border-emerald-200 shadow-sm ring-1 ring-emerald-200'
                            : 'bg-white border-gray-100 hover:border-gray-300'
                            }`}
                    >
                        <div className="flex items-center gap-3 mb-2">
                            <div className={`p-2 rounded-lg ${activeCategory === cat.id ? 'bg-emerald-600 text-white' : 'bg-gray-50 text-gray-400'}`}>
                                <cat.icon className="h-5 w-5" />
                            </div>
                            <span className={`font-bold ${activeCategory === cat.id ? 'text-emerald-900' : 'text-gray-700'}`}>
                                {cat.name}
                            </span>
                        </div>
                        <p className="text-xs text-gray-500">{cat.desc}</p>
                    </button>
                ))}
            </div>

            {/* Masters List */}
            <div className="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
                <div className="overflow-x-auto">
                    <table className="w-full text-left">
                        <thead className="bg-gray-50 border-b border-gray-200">
                            <tr>
                                <th className="px-6 py-4 text-xs font-bold text-gray-500 uppercase">Name</th>
                                <th className="px-6 py-4 text-xs font-bold text-gray-500 uppercase">Description</th>
                                <th className="px-6 py-4 text-xs font-bold text-gray-500 uppercase">Status</th>
                                <th className="px-6 py-4 text-xs font-bold text-gray-500 uppercase text-right">Actions</th>
                            </tr>
                        </thead>
                        <tbody className="divide-y divide-gray-100">
                            {loading ? (
                                <tr>
                                    <td colSpan={4} className="px-6 py-12 text-center text-gray-500">
                                        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-emerald-600 mx-auto"></div>
                                    </td>
                                </tr>
                            ) : filteredMasters.length === 0 ? (
                                <tr>
                                    <td colSpan={4} className="px-6 py-12 text-center text-gray-500">
                                        No entries found in this category.
                                    </td>
                                </tr>
                            ) : (
                                filteredMasters.map((master) => (
                                    <tr key={master.id} className="hover:bg-gray-50/50 transition-colors group">
                                        <td className="px-6 py-4">
                                            <div className="font-bold text-gray-900">{master.name}</div>
                                        </td>
                                        <td className="px-6 py-4">
                                            <div className="text-sm text-gray-500">{master.description || '-'}</div>
                                        </td>
                                        <td className="px-6 py-4">
                                            <span className={`px-2.5 py-1 text-[10px] font-bold rounded-full uppercase border ${master.is_active ? 'bg-emerald-50 text-emerald-700 border-emerald-100' : 'bg-gray-50 text-gray-600 border-gray-100'
                                                }`}>
                                                {master.is_active ? 'Active' : 'Inactive'}
                                            </span>
                                        </td>
                                        <td className="px-6 py-4 text-right">
                                            <div className="flex items-center justify-end gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                                                <button onClick={() => handleOpenEdit(master)} className="p-1.5 text-gray-400 hover:text-emerald-600 hover:bg-emerald-50 rounded-lg">
                                                    <Edit2 className="h-4 w-4" />
                                                </button>
                                                <button onClick={() => handleDelete(master.id)} className="p-1.5 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-lg">
                                                    <Trash2 className="h-4 w-4" />
                                                </button>
                                            </div>
                                        </td>
                                    </tr>
                                ))
                            )}
                        </tbody>
                    </table>
                </div>
            </div>

            {/* Modal */}
            {showModal && (
                <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40 backdrop-blur-sm animate-in fade-in duration-200">
                    <div className="bg-white rounded-2xl shadow-2xl max-w-lg w-full flex flex-col overflow-hidden">
                        <div className="px-6 py-4 border-b border-gray-100 flex items-center justify-between">
                            <h2 className="text-xl font-bold text-gray-900">{editingMaster ? 'Edit Entry' : 'Add New Entry'}</h2>
                            <button onClick={() => setShowModal(false)} className="p-2 hover:bg-gray-100 rounded-lg transition-colors">
                                <X className="h-5 w-5 text-gray-500" />
                            </button>
                        </div>

                        <form onSubmit={handleSubmit} className="p-6 space-y-4">
                            <div className="space-y-1.5">
                                <label className="text-xs font-bold text-gray-500 uppercase tracking-wider">Category</label>
                                <div className="px-4 py-2.5 bg-gray-100 border border-gray-200 rounded-xl text-gray-500 text-sm font-bold">
                                    {activeCategory}
                                </div>
                            </div>
                            <div className="space-y-1.5">
                                <label className="text-xs font-bold text-gray-500 uppercase tracking-wider">Name *</label>
                                <input
                                    type="text"
                                    required
                                    placeholder="Enter category name"
                                    className="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500 transition-all text-sm"
                                    value={formData.name}
                                    onChange={e => setFormData({ ...formData, name: e.target.value })}
                                />
                            </div>
                            <div className="space-y-1.5">
                                <label className="text-xs font-bold text-gray-500 uppercase tracking-wider">Description</label>
                                <textarea
                                    className="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500 transition-all text-sm min-h-[100px]"
                                    value={formData.description}
                                    onChange={e => setFormData({ ...formData, description: e.target.value })}
                                    placeholder="Optional description"
                                />
                            </div>
                            <div className="flex items-center gap-3 pt-2">
                                <input
                                    type="checkbox"
                                    id="is_active"
                                    className="h-4 w-4 text-emerald-600 focus:ring-emerald-500 border-gray-300 rounded"
                                    checked={formData.is_active}
                                    onChange={e => setFormData({ ...formData, is_active: e.target.checked })}
                                />
                                <label htmlFor="is_active" className="text-sm font-semibold text-gray-700">Active and available in forms</label>
                            </div>

                            <div className="pt-6 flex items-center justify-end gap-3">
                                <button
                                    type="button"
                                    onClick={() => setShowModal(false)}
                                    className="px-6 py-2.5 text-sm font-bold text-gray-700 hover:bg-gray-50 border border-transparent rounded-xl transition-all"
                                >
                                    Cancel
                                </button>
                                <button
                                    type="submit"
                                    className="px-8 py-2.5 bg-emerald-600 text-white text-sm font-bold rounded-xl hover:bg-emerald-700 transition-all shadow-lg shadow-emerald-600/20"
                                >
                                    {editingMaster ? 'Update' : 'Save'}
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            )}
        </div>
    );
};

export default NUIMasters;
