import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import '../../../../core/constants/app_colors.dart';
import '../cubit/registration_cubit.dart';

// ─────────────────────────────────────────────────────────────────────────────
// Data Models
// ─────────────────────────────────────────────────────────────────────────────

class _ModuleItem {
  final String id;
  final String name;
  final String description;
  final IconData icon;
  final bool isCore;

  const _ModuleItem({
    required this.id,
    required this.name,
    required this.description,
    required this.icon,
    this.isCore = false,
  });
}

class _AppItem {
  final String id;
  final String name;
  final String description;
  final IconData icon;
  final String platform;
  final IconData platformIcon;
  final List<String> requiredModules;
  final bool isCore;

  const _AppItem({
    required this.id,
    required this.name,
    required this.description,
    required this.icon,
    required this.platform,
    required this.platformIcon,
    this.requiredModules = const [],
    this.isCore = false,
  });
}

// ─────────────────────────────────────────────────────────────────────────────
// Static Data
// ─────────────────────────────────────────────────────────────────────────────

const _companyTypes = [
  'PT (Perseroan Terbatas)',
  'CV (Commanditaire Vennootschap)',
  'Firma',
  'Koperasi',
  'Yayasan',
  'Perorangan / UMKM',
  'Lainnya',
];

const _moduleGroups = <String, List<_ModuleItem>>{
  'Inti Sistem': [
    _ModuleItem(
      id: 'accounting',
      name: 'Akuntansi & Keuangan',
      description: 'Buku besar, jurnal, laporan keuangan, multi-mata uang',
      icon: Icons.account_balance_rounded,
      isCore: true,
    ),
  ],
  'Operasional': [
    _ModuleItem(
      id: 'sales',
      name: 'Penjualan & CRM',
      description: 'Lead, penawaran, pesanan penjualan, faktur, komisi',
      icon: Icons.trending_up_rounded,
    ),
    _ModuleItem(
      id: 'purchasing',
      name: 'Pembelian',
      description: 'PO, penerimaan barang, faktur pembelian, RFQ',
      icon: Icons.shopping_cart_rounded,
    ),
    _ModuleItem(
      id: 'stock',
      name: 'Inventori & Gudang',
      description: 'Stok, bin lokasi, opname, transfer, FIFO/FEFO',
      icon: Icons.warehouse_rounded,
    ),
    _ModuleItem(
      id: 'courier',
      name: 'Kurir & Pengiriman',
      description: 'Order pengiriman, tracking driver, POD',
      icon: Icons.local_shipping_rounded,
    ),
    _ModuleItem(
      id: 'pos',
      name: 'Point of Sale',
      description: 'Kasir ritel, struk fiskal, promo, printer termal',
      icon: Icons.point_of_sale_rounded,
    ),
  ],
  'SDM & Manajemen': [
    _ModuleItem(
      id: 'hrm',
      name: 'HR & Payroll',
      description: 'Karyawan, absensi, cuti, payroll, BPJS, PPh21',
      icon: Icons.people_rounded,
    ),
    _ModuleItem(
      id: 'fixed_assets',
      name: 'Aset Tetap',
      description: 'Master aset, penyusutan, revaluasi, pemeliharaan',
      icon: Icons.business_center_rounded,
    ),
    _ModuleItem(
      id: 'budgeting',
      name: 'Anggaran & Budget',
      description: 'Anggaran tahunan, realisasi, revisi, variance analysis',
      icon: Icons.pie_chart_rounded,
    ),
    _ModuleItem(
      id: 'project',
      name: 'Manajemen Proyek',
      description: 'RAB, milestone, OKR, KPI, approval workflow',
      icon: Icons.task_alt_rounded,
    ),
  ],
};

const _appGroups = <String, List<_AppItem>>{
  'Web': [
    _AppItem(
      id: 'web-ui',
      name: 'Admin Dashboard',
      description:
          'Panel administrasi utama — semua modul, laporan, master data, konfigurasi sistem',
      icon: Icons.dashboard_rounded,
      platform: 'Web (PWA)',
      platformIcon: Icons.language_rounded,
      isCore: true,
    ),
    _AppItem(
      id: 'app-management',
      name: 'Executive Dashboard',
      description:
          'Ringkasan eksekutif, grafik KPI, import data Excel/CSV, monitoring lintas departemen',
      icon: Icons.bar_chart_rounded,
      platform: 'Web + Mobile',
      platformIcon: Icons.devices_rounded,
    ),
  ],
  'Mobile — Penjualan': [
    _AppItem(
      id: 'app-sales-person',
      name: 'Sales Representative',
      description:
          'Lead, kunjungan customer, penawaran, target & realisasi, peta rute, komisi',
      icon: Icons.badge_rounded,
      platform: 'Android / iOS',
      platformIcon: Icons.smartphone_rounded,
      requiredModules: ['sales'],
    ),
    _AppItem(
      id: 'app-customer',
      name: 'Customer Portal',
      description:
          'Pelanggan lihat katalog, buat order, cek status invoice & pengiriman, riwayat transaksi',
      icon: Icons.people_alt_rounded,
      platform: 'Android / iOS',
      platformIcon: Icons.smartphone_rounded,
      requiredModules: ['sales'],
    ),
  ],
  'Mobile — Operasional': [
    _AppItem(
      id: 'app-pos',
      name: 'Point of Sale / Kasir',
      description:
          'Transaksi kasir ritel, printer termal, struk fiskal, manajemen promo & diskon',
      icon: Icons.point_of_sale_rounded,
      platform: 'Android / iOS',
      platformIcon: Icons.smartphone_rounded,
      requiredModules: ['pos'],
    ),
    _AppItem(
      id: 'app-opname',
      name: 'Stock Opname',
      description:
          'Hitung stok fisik, scan barcode, formulir opname per rak/kategori, rekonsiliasi',
      icon: Icons.inventory_2_rounded,
      platform: 'Android / iOS',
      platformIcon: Icons.smartphone_rounded,
      requiredModules: ['stock'],
    ),
    _AppItem(
      id: 'app-picking',
      name: 'Warehouse Picking',
      description:
          'Pick & pack pesanan, scan barcode, konfirmasi pengiriman — petugas gudang',
      icon: Icons.move_to_inbox_rounded,
      platform: 'Android / iOS',
      platformIcon: Icons.smartphone_rounded,
      requiredModules: ['stock'],
    ),
    _AppItem(
      id: 'app-courier',
      name: 'Delivery Driver',
      description:
          'Rute pengiriman, tracking GPS real-time, konfirmasi POD dengan foto & tanda tangan',
      icon: Icons.delivery_dining_rounded,
      platform: 'Android / iOS',
      platformIcon: Icons.smartphone_rounded,
      requiredModules: ['courier'],
    ),
  ],
  'Mobile — SDM & Pemasok': [
    _AppItem(
      id: 'app-employee',
      name: 'HR Self-Service',
      description:
          'Absensi dengan geofencing, pengajuan cuti, slip gaji, jadwal shift, chat tim',
      icon: Icons.person_pin_rounded,
      platform: 'Android / iOS',
      platformIcon: Icons.smartphone_rounded,
      requiredModules: ['hrm'],
    ),
    _AppItem(
      id: 'app-supplier',
      name: 'Supplier Portal',
      description:
          'Pemasok lihat PO, konfirmasi pengiriman, upload faktur, komunikasi tim pembelian',
      icon: Icons.storefront_rounded,
      platform: 'Android / iOS',
      platformIcon: Icons.smartphone_rounded,
      requiredModules: ['purchasing'],
    ),
  ],
};

// ─────────────────────────────────────────────────────────────────────────────
// Main Page
// ─────────────────────────────────────────────────────────────────────────────

class CreateClientPage extends StatefulWidget {
  const CreateClientPage({super.key});

  @override
  State<CreateClientPage> createState() => _CreateClientPageState();
}

class _CreateClientPageState extends State<CreateClientPage> {
  static const _totalSteps = 4;

  final _pageCtrl = PageController();
  int _step = 0;

  // Step 1 — Identitas
  final _formKey1 = GlobalKey<FormState>();
  final _nameCtrl = TextEditingController();
  final _codeCtrl = TextEditingController();
  final _npwpCtrl = TextEditingController();
  final _emailCtrl = TextEditingController();
  final _phoneCtrl = TextEditingController();
  final _addressCtrl = TextEditingController();
  final _websiteCtrl = TextEditingController();
  String _companyType = _companyTypes.first;
  bool _codeManuallyEdited = false;

  // Step 2 — Modul
  final Set<String> _selectedModules = {'accounting'};

  // Step 3 — Apps
  final Set<String> _selectedApps = {'web-ui'};

  bool _isSubmitting = false;

  @override
  void initState() {
    super.initState();
    _nameCtrl.addListener(_onNameChanged);
  }

  void _onNameChanged() {
    if (_codeManuallyEdited) return;
    final words = _nameCtrl.text
        .toUpperCase()
        .replaceAll(RegExp(r'[^A-Z0-9\s]'), '')
        .trim()
        .split(RegExp(r'\s+'))
        .where((w) => w.isNotEmpty)
        .toList();

    String code;
    if (words.isEmpty) {
      code = '';
    } else if (words.length == 1) {
      code = words[0].substring(0, words[0].length.clamp(0, 8));
    } else {
      code = words.take(3).map((w) => w.substring(0, w.length.clamp(0, 4))).join('-');
    }

    if (code != _codeCtrl.text) {
      _codeCtrl.value = TextEditingValue(
        text: code,
        selection: TextSelection.collapsed(offset: code.length),
      );
    }
  }

  @override
  void dispose() {
    _pageCtrl.dispose();
    _nameCtrl
      ..removeListener(_onNameChanged)
      ..dispose();
    _codeCtrl.dispose();
    _npwpCtrl.dispose();
    _emailCtrl.dispose();
    _phoneCtrl.dispose();
    _addressCtrl.dispose();
    _websiteCtrl.dispose();
    super.dispose();
  }

  void _goNext() {
    if (_step == 0 && !_formKey1.currentState!.validate()) return;
    if (_step < _totalSteps - 1) {
      setState(() => _step++);
      _pageCtrl.animateToPage(_step,
          duration: const Duration(milliseconds: 300), curve: Curves.easeInOut);
    }
  }

  void _goBack() {
    if (_step > 0) {
      setState(() => _step--);
      _pageCtrl.animateToPage(_step,
          duration: const Duration(milliseconds: 300), curve: Curves.easeInOut);
    } else {
      Navigator.pop(context);
    }
  }

  Future<void> _submit() async {
    setState(() => _isSubmitting = true);
    context.read<RegistrationCubit>().createClient(
          code: _codeCtrl.text.trim(),
          name: _nameCtrl.text.trim(),
          companyType: _companyType,
          npwp: _npwpCtrl.text.trim().isEmpty ? null : _npwpCtrl.text.trim(),
          email: _emailCtrl.text.trim().isEmpty ? null : _emailCtrl.text.trim(),
          phone: _phoneCtrl.text.trim().isEmpty ? null : _phoneCtrl.text.trim(),
          address: _addressCtrl.text.trim().isEmpty ? null : _addressCtrl.text.trim(),
          website: _websiteCtrl.text.trim().isEmpty ? null : _websiteCtrl.text.trim(),
          modules: _selectedModules.toList(),
          apps: _selectedApps.toList(),
        );
    if (mounted) Navigator.pop(context);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        leading: IconButton(
          icon: const Icon(Icons.close_rounded),
          onPressed: () => Navigator.pop(context),
          tooltip: 'Tutup',
        ),
        title: const Text('Buat Client Baru'),
        centerTitle: false,
        bottom: PreferredSize(
          preferredSize: const Size.fromHeight(3),
          child: _ProgressBar(step: _step, total: _totalSteps),
        ),
      ),
      body: Column(
        children: [
          _StepHeader(step: _step, total: _totalSteps),
          const Divider(height: 1),
          Expanded(
            child: PageView(
              controller: _pageCtrl,
              physics: const NeverScrollableScrollPhysics(),
              children: [
                _Step1Identity(
                  formKey: _formKey1,
                  nameCtrl: _nameCtrl,
                  codeCtrl: _codeCtrl,
                  npwpCtrl: _npwpCtrl,
                  emailCtrl: _emailCtrl,
                  phoneCtrl: _phoneCtrl,
                  addressCtrl: _addressCtrl,
                  websiteCtrl: _websiteCtrl,
                  companyType: _companyType,
                  onCompanyTypeChanged: (v) => setState(() => _companyType = v),
                  onCodeEdited: () => setState(() => _codeManuallyEdited = true),
                ),
                _Step2Modules(
                  selected: _selectedModules,
                  onToggle: (id, isCore) {
                    if (isCore) return;
                    setState(() {
                      if (_selectedModules.contains(id)) {
                        _selectedModules.remove(id);
                        // Hapus app yang bergantung pada modul ini jika tidak ada modul lain
                        _cleanDependentApps(id);
                      } else {
                        _selectedModules.add(id);
                      }
                    });
                  },
                ),
                _Step3Apps(
                  selectedApps: _selectedApps,
                  selectedModules: _selectedModules,
                  onToggle: (id, isCore) {
                    if (isCore) return;
                    setState(() {
                      if (_selectedApps.contains(id)) {
                        _selectedApps.remove(id);
                      } else {
                        _selectedApps.add(id);
                      }
                    });
                  },
                ),
                _Step4Review(
                  name: _nameCtrl.text,
                  code: _codeCtrl.text,
                  companyType: _companyType,
                  npwp: _npwpCtrl.text,
                  email: _emailCtrl.text,
                  phone: _phoneCtrl.text,
                  address: _addressCtrl.text,
                  website: _websiteCtrl.text,
                  modules: _selectedModules,
                  apps: _selectedApps,
                ),
              ],
            ),
          ),
        ],
      ),
      bottomNavigationBar: SafeArea(
        child: Padding(
          padding: const EdgeInsets.fromLTRB(20, 12, 20, 16),
          child: Row(
            children: [
              OutlinedButton(
                onPressed: _isSubmitting ? null : _goBack,
                style: OutlinedButton.styleFrom(
                  padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 14),
                  side: const BorderSide(color: AppColors.neutral300),
                  foregroundColor: AppColors.neutral700,
                  shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(12)),
                ),
                child: Text(_step == 0 ? 'Batal' : 'Kembali'),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: _step < _totalSteps - 1
                    ? ElevatedButton(
                        onPressed: _goNext,
                        style: ElevatedButton.styleFrom(
                          backgroundColor: AppColors.primary700,
                          foregroundColor: AppColors.white,
                          padding: const EdgeInsets.symmetric(vertical: 14),
                          shape: RoundedRectangleBorder(
                              borderRadius: BorderRadius.circular(12)),
                        ),
                        child: Row(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            Text(_nextLabel()),
                            const SizedBox(width: 6),
                            const Icon(Icons.arrow_forward_rounded, size: 18),
                          ],
                        ),
                      )
                    : ElevatedButton(
                        onPressed: _isSubmitting ? null : _submit,
                        style: ElevatedButton.styleFrom(
                          backgroundColor: AppColors.success,
                          foregroundColor: AppColors.white,
                          padding: const EdgeInsets.symmetric(vertical: 14),
                          shape: RoundedRectangleBorder(
                              borderRadius: BorderRadius.circular(12)),
                        ),
                        child: _isSubmitting
                            ? const SizedBox(
                                width: 20,
                                height: 20,
                                child: CircularProgressIndicator(
                                    strokeWidth: 2, color: AppColors.white),
                              )
                            : const Row(
                                mainAxisAlignment: MainAxisAlignment.center,
                                children: [
                                  Icon(Icons.check_rounded, size: 18),
                                  SizedBox(width: 6),
                                  Text('Buat Client'),
                                ],
                              ),
                      ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  String _nextLabel() {
    return switch (_step) {
      0 => 'Pilih Modul',
      1 => 'Pilih Aplikasi',
      _ => 'Review',
    };
  }

  void _cleanDependentApps(String removedModule) {
    final allApps = _appGroups.values.expand((list) => list);
    for (final app in allApps) {
      if (app.isCore) continue;
      if (app.requiredModules.contains(removedModule)) {
        final otherModulesSatisfied =
            app.requiredModules.where((m) => m != removedModule).every(_selectedModules.contains);
        if (!otherModulesSatisfied) {
          _selectedApps.remove(app.id);
        }
      }
    }
  }
}

// ─────────────────────────────────────────────────────────────────────────────
// Step Header
// ─────────────────────────────────────────────────────────────────────────────

class _StepHeader extends StatelessWidget {
  final int step;
  final int total;

  const _StepHeader({required this.step, required this.total});

  static const _labels = ['Identitas', 'Modul', 'Aplikasi', 'Konfirmasi'];

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 16),
      child: Row(
        children: List.generate(total * 2 - 1, (i) {
          if (i.isOdd) {
            return Expanded(
              child: AnimatedContainer(
                duration: const Duration(milliseconds: 300),
                height: 2,
                color: i ~/ 2 < step ? AppColors.primary700 : AppColors.neutral200,
              ),
            );
          }
          final idx = i ~/ 2;
          final isDone = idx < step;
          final isActive = idx == step;
          return Column(
            children: [
              AnimatedContainer(
                duration: const Duration(milliseconds: 250),
                width: isActive ? 32 : 24,
                height: isActive ? 32 : 24,
                decoration: BoxDecoration(
                  color: isDone || isActive ? AppColors.primary700 : AppColors.neutral200,
                  shape: BoxShape.circle,
                ),
                child: Center(
                  child: isDone
                      ? const Icon(Icons.check_rounded, size: 14, color: AppColors.white)
                      : Text(
                          '${idx + 1}',
                          style: TextStyle(
                            fontSize: isActive ? 13 : 11,
                            fontWeight: FontWeight.w700,
                            color: isActive ? AppColors.white : AppColors.neutral500,
                          ),
                        ),
                ),
              ),
              const SizedBox(height: 4),
              Text(
                _labels[idx],
                style: TextStyle(
                  fontSize: 10,
                  fontWeight: isActive ? FontWeight.w700 : FontWeight.normal,
                  color: isActive ? AppColors.primary700 : AppColors.neutral400,
                ),
              ),
            ],
          );
        }),
      ),
    );
  }
}

// ─────────────────────────────────────────────────────────────────────────────
// Step 1: Identitas Perusahaan
// ─────────────────────────────────────────────────────────────────────────────

class _Step1Identity extends StatelessWidget {
  final GlobalKey<FormState> formKey;
  final TextEditingController nameCtrl;
  final TextEditingController codeCtrl;
  final TextEditingController npwpCtrl;
  final TextEditingController emailCtrl;
  final TextEditingController phoneCtrl;
  final TextEditingController addressCtrl;
  final TextEditingController websiteCtrl;
  final String companyType;
  final ValueChanged<String> onCompanyTypeChanged;
  final VoidCallback onCodeEdited;

  const _Step1Identity({
    required this.formKey,
    required this.nameCtrl,
    required this.codeCtrl,
    required this.npwpCtrl,
    required this.emailCtrl,
    required this.phoneCtrl,
    required this.addressCtrl,
    required this.websiteCtrl,
    required this.companyType,
    required this.onCompanyTypeChanged,
    required this.onCodeEdited,
  });

  @override
  Widget build(BuildContext context) {
    return Form(
      key: formKey,
      child: ListView(
        padding: const EdgeInsets.fromLTRB(20, 20, 20, 24),
        children: [
          _SectionTitle(
            icon: Icons.business_rounded,
            title: 'Identitas Perusahaan',
            subtitle: 'Data resmi perusahaan client',
          ),
          const SizedBox(height: 20),

          // Tipe perusahaan
          DropdownButtonFormField<String>(
            value: companyType,
            onChanged: (v) => onCompanyTypeChanged(v!),
            decoration: const InputDecoration(
              labelText: 'Tipe Perusahaan *',
              prefixIcon: Icon(Icons.domain_rounded),
            ),
            items: _companyTypes
                .map((t) => DropdownMenuItem(value: t, child: Text(t, style: const TextStyle(fontSize: 14))))
                .toList(),
          ),
          const SizedBox(height: 14),

          // Nama Perusahaan
          TextFormField(
            controller: nameCtrl,
            textCapitalization: TextCapitalization.words,
            decoration: const InputDecoration(
              labelText: 'Nama Perusahaan *',
              hintText: 'misal: PT Maju Bersama',
              prefixIcon: Icon(Icons.apartment_rounded),
              helperText: 'Kode akan digenerate otomatis',
            ),
            validator: (v) => v == null || v.trim().isEmpty ? 'Nama wajib diisi' : null,
          ),
          const SizedBox(height: 14),

          // Kode Perusahaan
          TextFormField(
            controller: codeCtrl,
            textCapitalization: TextCapitalization.characters,
            onChanged: (_) => onCodeEdited(),
            decoration: const InputDecoration(
              labelText: 'Kode Perusahaan *',
              hintText: 'misal: PT-MAJU',
              prefixIcon: Icon(Icons.tag_rounded),
              helperText: 'Tidak dapat diubah setelah dibuat',
            ),
            validator: (v) {
              if (v == null || v.trim().isEmpty) return 'Kode wajib diisi';
              if (v.trim().length < 2) return 'Minimal 2 karakter';
              return null;
            },
          ),
          const SizedBox(height: 20),

          _FieldDivider(label: 'Legalitas'),
          const SizedBox(height: 14),

          TextFormField(
            controller: npwpCtrl,
            keyboardType: TextInputType.number,
            decoration: const InputDecoration(
              labelText: 'NPWP',
              hintText: '00.000.000.0-000.000',
              prefixIcon: Icon(Icons.receipt_long_rounded),
            ),
          ),
          const SizedBox(height: 20),

          _FieldDivider(label: 'Kontak'),
          const SizedBox(height: 14),

          TextFormField(
            controller: emailCtrl,
            keyboardType: TextInputType.emailAddress,
            decoration: const InputDecoration(
              labelText: 'Email Perusahaan',
              hintText: 'info@perusahaan.com',
              prefixIcon: Icon(Icons.email_outlined),
            ),
            validator: (v) {
              if (v == null || v.trim().isEmpty) return null;
              if (!v.contains('@')) return 'Format email tidak valid';
              return null;
            },
          ),
          const SizedBox(height: 14),

          TextFormField(
            controller: phoneCtrl,
            keyboardType: TextInputType.phone,
            decoration: const InputDecoration(
              labelText: 'Nomor Telepon',
              hintText: '021-xxxx-xxxx',
              prefixIcon: Icon(Icons.phone_outlined),
            ),
          ),
          const SizedBox(height: 14),

          TextFormField(
            controller: websiteCtrl,
            keyboardType: TextInputType.url,
            decoration: const InputDecoration(
              labelText: 'Website',
              hintText: 'https://perusahaan.com',
              prefixIcon: Icon(Icons.language_rounded),
            ),
          ),
          const SizedBox(height: 14),

          TextFormField(
            controller: addressCtrl,
            maxLines: 3,
            textCapitalization: TextCapitalization.sentences,
            decoration: const InputDecoration(
              labelText: 'Alamat Perusahaan',
              hintText: 'Jl. ...',
              prefixIcon: Icon(Icons.location_on_outlined),
              alignLabelWithHint: true,
            ),
          ),
        ],
      ),
    );
  }
}

// ─────────────────────────────────────────────────────────────────────────────
// Step 2: Modul & Fitur
// ─────────────────────────────────────────────────────────────────────────────

class _Step2Modules extends StatelessWidget {
  final Set<String> selected;
  final void Function(String id, bool isCore) onToggle;

  const _Step2Modules({required this.selected, required this.onToggle});

  @override
  Widget build(BuildContext context) {
    return ListView(
      padding: const EdgeInsets.fromLTRB(20, 20, 20, 24),
      children: [
        _SectionTitle(
          icon: Icons.extension_rounded,
          title: 'Modul & Fitur',
          subtitle: '${selected.length} modul aktif — pilih sesuai paket client',
        ),
        const SizedBox(height: 20),
        ..._moduleGroups.entries.map((group) => Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                _GroupLabel(group.key),
                const SizedBox(height: 10),
                ...group.value.map((mod) => _ModuleCard(
                      item: mod,
                      isSelected: selected.contains(mod.id),
                      onTap: () => onToggle(mod.id, mod.isCore),
                    )),
                const SizedBox(height: 18),
              ],
            )),
      ],
    );
  }
}

class _ModuleCard extends StatelessWidget {
  final _ModuleItem item;
  final bool isSelected;
  final VoidCallback onTap;

  const _ModuleCard({required this.item, required this.isSelected, required this.onTap});

  @override
  Widget build(BuildContext context) {
    final isCore = item.isCore;
    return GestureDetector(
      onTap: isCore ? null : onTap,
      child: AnimatedContainer(
        duration: const Duration(milliseconds: 180),
        margin: const EdgeInsets.only(bottom: 10),
        padding: const EdgeInsets.all(14),
        decoration: BoxDecoration(
          color: isSelected ? AppColors.primary50 : (isCore ? AppColors.neutral100 : AppColors.surface),
          borderRadius: BorderRadius.circular(12),
          border: Border.all(
              color: isSelected ? AppColors.primary700 : AppColors.neutral200,
              width: isSelected ? 1.5 : 1),
        ),
        child: Row(
          children: [
            Container(
              width: 42,
              height: 42,
              decoration: BoxDecoration(
                color: isSelected ? AppColors.primary700 : AppColors.neutral200,
                borderRadius: BorderRadius.circular(10),
              ),
              child: Icon(item.icon,
                  color: isSelected ? AppColors.white : AppColors.neutral500, size: 22),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      Expanded(
                        child: Text(item.name,
                            style: TextStyle(
                              fontSize: 14,
                              fontWeight: FontWeight.w600,
                              color: isSelected ? AppColors.primary700 : AppColors.neutral900,
                            )),
                      ),
                      if (isCore)
                        _Badge(label: 'WAJIB', color: AppColors.primary700, bg: AppColors.primary100),
                    ],
                  ),
                  const SizedBox(height: 2),
                  Text(item.description,
                      style: const TextStyle(fontSize: 12, color: AppColors.neutral500)),
                ],
              ),
            ),
            const SizedBox(width: 8),
            if (isCore)
              const Icon(Icons.lock_rounded, size: 18, color: AppColors.neutral400)
            else
              _Checkbox(checked: isSelected),
          ],
        ),
      ),
    );
  }
}

// ─────────────────────────────────────────────────────────────────────────────
// Step 3: Aplikasi & Interface
// ─────────────────────────────────────────────────────────────────────────────

class _Step3Apps extends StatelessWidget {
  final Set<String> selectedApps;
  final Set<String> selectedModules;
  final void Function(String id, bool isCore) onToggle;

  const _Step3Apps({
    required this.selectedApps,
    required this.selectedModules,
    required this.onToggle,
  });

  @override
  Widget build(BuildContext context) {
    final enabledCount = selectedApps.length;
    return ListView(
      padding: const EdgeInsets.fromLTRB(20, 20, 20, 24),
      children: [
        _SectionTitle(
          icon: Icons.apps_rounded,
          title: 'Aplikasi & Interface',
          subtitle: '$enabledCount aplikasi diaktifkan — pilih yang termasuk dalam paket',
        ),
        const SizedBox(height: 12),
        // Info hint modul
        Container(
          padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
          decoration: BoxDecoration(
            color: AppColors.infoLight,
            borderRadius: BorderRadius.circular(8),
            border: Border.all(color: AppColors.info.withValues(alpha: 0.25)),
          ),
          child: const Row(
            children: [
              Icon(Icons.info_outline_rounded, size: 16, color: AppColors.info),
              SizedBox(width: 8),
              Expanded(
                child: Text(
                  'Aplikasi yang memerlukan modul tertentu hanya bisa diaktifkan jika modul tersebut dipilih di langkah sebelumnya.',
                  style: TextStyle(fontSize: 12, color: AppColors.infoBase),
                ),
              ),
            ],
          ),
        ),
        const SizedBox(height: 20),
        ..._appGroups.entries.map((group) => Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                _GroupLabel(group.key),
                const SizedBox(height: 10),
                ...group.value.map((app) {
                  final moduleSatisfied = app.requiredModules.isEmpty ||
                      app.requiredModules.every(selectedModules.contains);
                  return _AppCard(
                    item: app,
                    isSelected: selectedApps.contains(app.id),
                    moduleSatisfied: moduleSatisfied,
                    onTap: () => moduleSatisfied ? onToggle(app.id, app.isCore) : null,
                  );
                }),
                const SizedBox(height: 18),
              ],
            )),
      ],
    );
  }
}

class _AppCard extends StatelessWidget {
  final _AppItem item;
  final bool isSelected;
  final bool moduleSatisfied;
  final VoidCallback? onTap;

  const _AppCard({
    required this.item,
    required this.isSelected,
    required this.moduleSatisfied,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    final isCore = item.isCore;
    final locked = !moduleSatisfied;

    Color borderColor;
    Color bgColor;
    if (locked) {
      borderColor = AppColors.neutral200;
      bgColor = AppColors.neutral100;
    } else if (isSelected) {
      borderColor = AppColors.primary700;
      bgColor = AppColors.primary50;
    } else {
      borderColor = AppColors.neutral200;
      bgColor = AppColors.surface;
    }

    return GestureDetector(
      onTap: locked ? null : onTap,
      child: AnimatedContainer(
        duration: const Duration(milliseconds: 180),
        margin: const EdgeInsets.only(bottom: 10),
        padding: const EdgeInsets.all(14),
        decoration: BoxDecoration(
          color: bgColor,
          borderRadius: BorderRadius.circular(12),
          border: Border.all(
              color: borderColor, width: isSelected && !locked ? 1.5 : 1),
        ),
        child: Row(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // App icon
            Container(
              width: 44,
              height: 44,
              decoration: BoxDecoration(
                color: locked
                    ? AppColors.neutral200
                    : (isSelected ? AppColors.primary700 : AppColors.neutral200),
                borderRadius: BorderRadius.circular(10),
              ),
              child: Icon(item.icon,
                  color: locked
                      ? AppColors.neutral400
                      : (isSelected ? AppColors.white : AppColors.neutral500),
                  size: 22),
            ),
            const SizedBox(width: 12),

            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // Baris nama + badge
                  Row(
                    children: [
                      Expanded(
                        child: Text(
                          item.name,
                          style: TextStyle(
                            fontSize: 14,
                            fontWeight: FontWeight.w600,
                            color: locked
                                ? AppColors.neutral400
                                : (isSelected ? AppColors.primary700 : AppColors.neutral900),
                          ),
                        ),
                      ),
                      if (isCore)
                        _Badge(label: 'WAJIB', color: AppColors.primary700, bg: AppColors.primary100)
                      else if (locked)
                        _Badge(
                          label: 'PERLU MODUL',
                          color: AppColors.neutral500,
                          bg: AppColors.neutral200,
                        ),
                    ],
                  ),
                  const SizedBox(height: 3),
                  Text(item.description,
                      style: TextStyle(
                          fontSize: 12,
                          color: locked ? AppColors.neutral400 : AppColors.neutral500)),
                  const SizedBox(height: 6),
                  // Platform badge
                  Row(
                    children: [
                      Icon(item.platformIcon,
                          size: 13,
                          color: locked ? AppColors.neutral400 : AppColors.neutral400),
                      const SizedBox(width: 4),
                      Text(
                        item.platform,
                        style: TextStyle(
                          fontSize: 11,
                          fontWeight: FontWeight.w500,
                          color: locked ? AppColors.neutral400 : AppColors.neutral500,
                        ),
                      ),
                    ],
                  ),
                  // Modul yang diperlukan jika locked
                  if (locked && item.requiredModules.isNotEmpty) ...[
                    const SizedBox(height: 5),
                    _ModuleRequiredHint(item.requiredModules),
                  ],
                ],
              ),
            ),

            const SizedBox(width: 8),
            if (isCore)
              const Icon(Icons.lock_rounded, size: 18, color: AppColors.neutral400)
            else if (locked)
              const Icon(Icons.lock_outline_rounded, size: 18, color: AppColors.neutral300)
            else
              _Checkbox(checked: isSelected),
          ],
        ),
      ),
    );
  }
}

class _ModuleRequiredHint extends StatelessWidget {
  final List<String> requiredModules;

  const _ModuleRequiredHint(this.requiredModules);

  static const _moduleNames = <String, String>{
    'sales': 'Penjualan & CRM',
    'purchasing': 'Pembelian',
    'stock': 'Inventori & Gudang',
    'hrm': 'HR & Payroll',
    'courier': 'Kurir & Pengiriman',
    'pos': 'Point of Sale',
    'fixed_assets': 'Aset Tetap',
    'budgeting': 'Anggaran',
    'project': 'Proyek',
  };

  @override
  Widget build(BuildContext context) {
    final names = requiredModules
        .map((id) => _moduleNames[id] ?? id)
        .join(', ');
    return Text(
      'Aktifkan modul: $names',
      style: const TextStyle(
        fontSize: 11,
        fontStyle: FontStyle.italic,
        color: AppColors.neutral400,
      ),
    );
  }
}

// ─────────────────────────────────────────────────────────────────────────────
// Step 4: Review & Konfirmasi
// ─────────────────────────────────────────────────────────────────────────────

class _Step4Review extends StatelessWidget {
  final String name;
  final String code;
  final String companyType;
  final String npwp;
  final String email;
  final String phone;
  final String address;
  final String website;
  final Set<String> modules;
  final Set<String> apps;

  const _Step4Review({
    required this.name,
    required this.code,
    required this.companyType,
    required this.npwp,
    required this.email,
    required this.phone,
    required this.address,
    required this.website,
    required this.modules,
    required this.apps,
  });

  @override
  Widget build(BuildContext context) {
    final allModules = _moduleGroups.values.expand((l) => l).toList();
    final allApps = _appGroups.values.expand((l) => l).toList();
    final selectedModuleItems = allModules.where((m) => modules.contains(m.id)).toList();
    final selectedAppItems = allApps.where((a) => apps.contains(a.id)).toList();

    return ListView(
      padding: const EdgeInsets.fromLTRB(20, 20, 20, 24),
      children: [
        _SectionTitle(
          icon: Icons.checklist_rounded,
          title: 'Review & Konfirmasi',
          subtitle: 'Periksa kembali semua data sebelum membuat client',
        ),
        const SizedBox(height: 20),

        // ── Identitas ─────────────────────────────────────────────────────
        _ReviewCard(
          title: 'Identitas Perusahaan',
          icon: Icons.business_rounded,
          children: [
            _ReviewRow(label: 'Tipe', value: companyType),
            _ReviewRow(label: 'Nama', value: name),
            _ReviewRow(label: 'Kode', value: code, mono: true),
            if (npwp.isNotEmpty) _ReviewRow(label: 'NPWP', value: npwp),
            if (email.isNotEmpty) _ReviewRow(label: 'Email', value: email),
            if (phone.isNotEmpty) _ReviewRow(label: 'Telepon', value: phone),
            if (website.isNotEmpty) _ReviewRow(label: 'Website', value: website),
            if (address.isNotEmpty) _ReviewRow(label: 'Alamat', value: address),
          ],
        ),
        const SizedBox(height: 14),

        // ── Modul ─────────────────────────────────────────────────────────
        _ReviewCard(
          title: '${selectedModuleItems.length} Modul Diaktifkan',
          icon: Icons.extension_rounded,
          children: [
            Wrap(
              spacing: 8,
              runSpacing: 8,
              children: selectedModuleItems
                  .map((m) => _ItemChip(label: m.name, icon: m.icon))
                  .toList(),
            ),
          ],
        ),
        const SizedBox(height: 14),

        // ── Aplikasi ──────────────────────────────────────────────────────
        _ReviewCard(
          title: '${selectedAppItems.length} Aplikasi Diaktifkan',
          icon: Icons.apps_rounded,
          children: [
            ...selectedAppItems.map((a) => Padding(
                  padding: const EdgeInsets.only(bottom: 10),
                  child: Row(
                    children: [
                      Container(
                        width: 32,
                        height: 32,
                        decoration: BoxDecoration(
                          color: AppColors.primary100,
                          borderRadius: BorderRadius.circular(8),
                        ),
                        child: Icon(a.icon, size: 16, color: AppColors.primary700),
                      ),
                      const SizedBox(width: 10),
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(a.name,
                                style: const TextStyle(
                                    fontSize: 13,
                                    fontWeight: FontWeight.w600,
                                    color: AppColors.neutral900)),
                            Row(
                              children: [
                                Icon(a.platformIcon, size: 11, color: AppColors.neutral400),
                                const SizedBox(width: 3),
                                Text(a.platform,
                                    style: const TextStyle(
                                        fontSize: 11, color: AppColors.neutral500)),
                              ],
                            ),
                          ],
                        ),
                      ),
                    ],
                  ),
                )),
          ],
        ),
        const SizedBox(height: 20),

        // Note
        Container(
          padding: const EdgeInsets.all(14),
          decoration: BoxDecoration(
            color: AppColors.warningLight,
            borderRadius: BorderRadius.circular(10),
            border: Border.all(color: AppColors.warning.withValues(alpha: 0.3)),
          ),
          child: const Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Icon(Icons.warning_amber_rounded, size: 18, color: AppColors.warningBase),
              SizedBox(width: 8),
              Expanded(
                child: Text(
                  'Kode perusahaan tidak dapat diubah setelah dibuat. Modul dan aplikasi dapat disesuaikan dari panel admin setelah client aktif.',
                  style: TextStyle(fontSize: 13, color: AppColors.warningBase),
                ),
              ),
            ],
          ),
        ),
      ],
    );
  }
}

// ─────────────────────────────────────────────────────────────────────────────
// Shared UI Components
// ─────────────────────────────────────────────────────────────────────────────

class _ProgressBar extends StatelessWidget {
  final int step;
  final int total;

  const _ProgressBar({required this.step, required this.total});

  @override
  Widget build(BuildContext context) {
    return LinearProgressIndicator(
      value: (step + 1) / total,
      backgroundColor: AppColors.primary100,
      color: AppColors.primary700,
      minHeight: 3,
    );
  }
}

class _SectionTitle extends StatelessWidget {
  final IconData icon;
  final String title;
  final String subtitle;

  const _SectionTitle({
    required this.icon,
    required this.title,
    required this.subtitle,
  });

  @override
  Widget build(BuildContext context) {
    return Row(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Container(
          width: 40,
          height: 40,
          decoration: BoxDecoration(
            color: AppColors.primary100,
            borderRadius: BorderRadius.circular(10),
          ),
          child: Icon(icon, color: AppColors.primary700, size: 20),
        ),
        const SizedBox(width: 12),
        Expanded(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(title,
                  style: const TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.w700,
                      color: AppColors.neutral900)),
              const SizedBox(height: 2),
              Text(subtitle,
                  style: const TextStyle(fontSize: 13, color: AppColors.neutral500)),
            ],
          ),
        ),
      ],
    );
  }
}

class _GroupLabel extends StatelessWidget {
  final String text;

  const _GroupLabel(this.text);

  @override
  Widget build(BuildContext context) {
    return Text(
      text.toUpperCase(),
      style: const TextStyle(
        fontSize: 10,
        fontWeight: FontWeight.w700,
        color: AppColors.neutral400,
        letterSpacing: 1.2,
      ),
    );
  }
}

class _FieldDivider extends StatelessWidget {
  final String label;

  const _FieldDivider({required this.label});

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        const Expanded(child: Divider()),
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 12),
          child: Text(
            label,
            style: const TextStyle(
                fontSize: 11, color: AppColors.neutral400, fontWeight: FontWeight.w500),
          ),
        ),
        const Expanded(child: Divider()),
      ],
    );
  }
}

class _Badge extends StatelessWidget {
  final String label;
  final Color color;
  final Color bg;

  const _Badge({required this.label, required this.color, required this.bg});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
      decoration: BoxDecoration(color: bg, borderRadius: BorderRadius.circular(4)),
      child: Text(
        label,
        style: TextStyle(
            fontSize: 9, fontWeight: FontWeight.w700, color: color, letterSpacing: 0.4),
      ),
    );
  }
}

class _Checkbox extends StatelessWidget {
  final bool checked;

  const _Checkbox({required this.checked});

  @override
  Widget build(BuildContext context) {
    return AnimatedContainer(
      duration: const Duration(milliseconds: 180),
      width: 22,
      height: 22,
      decoration: BoxDecoration(
        color: checked ? AppColors.primary700 : AppColors.surface,
        borderRadius: BorderRadius.circular(6),
        border: Border.all(
          color: checked ? AppColors.primary700 : AppColors.neutral300,
          width: 1.5,
        ),
      ),
      child: checked
          ? const Icon(Icons.check_rounded, size: 14, color: AppColors.white)
          : null,
    );
  }
}

class _ReviewCard extends StatelessWidget {
  final String title;
  final IconData icon;
  final List<Widget> children;

  const _ReviewCard({required this.title, required this.icon, required this.children});

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
        color: AppColors.surface,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: AppColors.neutral200),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Padding(
            padding: const EdgeInsets.fromLTRB(16, 12, 16, 10),
            child: Row(
              children: [
                Icon(icon, size: 15, color: AppColors.primary700),
                const SizedBox(width: 8),
                Text(title,
                    style: const TextStyle(
                        fontSize: 13,
                        fontWeight: FontWeight.w700,
                        color: AppColors.primary700)),
              ],
            ),
          ),
          const Divider(height: 1),
          Padding(
            padding: const EdgeInsets.fromLTRB(16, 12, 16, 14),
            child: Column(
                crossAxisAlignment: CrossAxisAlignment.start, children: children),
          ),
        ],
      ),
    );
  }
}

class _ReviewRow extends StatelessWidget {
  final String label;
  final String value;
  final bool mono;

  const _ReviewRow({required this.label, required this.value, this.mono = false});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 8),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          SizedBox(
            width: 72,
            child: Text(label,
                style: const TextStyle(fontSize: 12, color: AppColors.neutral500)),
          ),
          Expanded(
            child: Text(
              value,
              style: TextStyle(
                fontSize: 13,
                fontWeight: FontWeight.w500,
                color: AppColors.neutral900,
                fontFamily: mono ? 'monospace' : null,
              ),
            ),
          ),
        ],
      ),
    );
  }
}

class _ItemChip extends StatelessWidget {
  final String label;
  final IconData icon;

  const _ItemChip({required this.label, required this.icon});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 6),
      decoration: BoxDecoration(
        color: AppColors.primary50,
        borderRadius: BorderRadius.circular(20),
        border: Border.all(color: AppColors.primary100),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(icon, size: 13, color: AppColors.primary700),
          const SizedBox(width: 5),
          Text(label,
              style: const TextStyle(
                  fontSize: 12, fontWeight: FontWeight.w500, color: AppColors.primary700)),
        ],
      ),
    );
  }
}
