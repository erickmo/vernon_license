import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import '../../../../core/constants/app_colors.dart';
import '../../domain/entities/product_entity.dart';
import '../cubit/product_cubit.dart';

// ── Page ────────────────────────────────────────────────────────────────────

class ProductFormPage extends StatefulWidget {
  final ProductEntity? initialProduct;
  const ProductFormPage({super.key, this.initialProduct});

  @override
  State<ProductFormPage> createState() => _ProductFormPageState();
}

class _ProductFormPageState extends State<ProductFormPage> {
  final _formKey = GlobalKey<FormState>();
  bool _isSaving = false;
  bool _isEdit = false;
  bool _slugEdited = false;

  // ── Basic info ─────────────────────────────────────────────────────────
  late final TextEditingController _nameCtrl;
  late final TextEditingController _slugCtrl;
  late final TextEditingController _descCtrl;
  bool _isActive = true;

  // ── Plans ──────────────────────────────────────────────────────────────
  List<String> _plans = [];
  final _planInputCtrl = TextEditingController();

  // ── Pricing — parallel maps keyed by plan name ─────────────────────────
  final Map<String, TextEditingController> _priceBase = {};
  final Map<String, TextEditingController> _pricePerUser = {};

  // ── Modules — parallel flat lists ──────────────────────────────────────
  final List<TextEditingController> _modKeys = [];
  final List<TextEditingController> _modNames = [];
  final List<TextEditingController> _modDescs = [];

  // ── Apps — parallel flat lists ─────────────────────────────────────────
  final List<TextEditingController> _appKeys = [];
  final List<TextEditingController> _appNames = [];
  final List<TextEditingController> _appDescs = [];

  // ── Init ───────────────────────────────────────────────────────────────

  @override
  void initState() {
    super.initState();
    _isEdit = widget.initialProduct != null;
    final p = widget.initialProduct;

    _nameCtrl = TextEditingController(text: p?.name ?? '');
    _slugCtrl = TextEditingController(text: p?.slug ?? '');
    _descCtrl = TextEditingController(text: p?.description ?? '');
    _isActive = p?.isActive ?? true;
    _slugEdited = _isEdit;

    if (p != null) {
      _plans = List.from(p.availablePlans);

      for (final plan in _plans) {
        final pr = p.basePricing[plan];
        _priceBase[plan] = TextEditingController(
          text: pr?.basePrice.toStringAsFixed(0) ?? '0',
        );
        _pricePerUser[plan] = TextEditingController(
          text: pr?.perUserPrice.toStringAsFixed(0) ?? '0',
        );
      }

      for (final m in p.availableModules) {
        _modKeys.add(TextEditingController(text: m.key));
        _modNames.add(TextEditingController(text: m.name));
        _modDescs.add(TextEditingController(text: m.description));
      }

      for (final a in p.availableApps) {
        _appKeys.add(TextEditingController(text: a.key));
        _appNames.add(TextEditingController(text: a.name));
        _appDescs.add(TextEditingController(text: a.description));
      }
    }

    _nameCtrl.addListener(_autoSlug);
  }

  // ── Helpers ────────────────────────────────────────────────────────────

  void _autoSlug() {
    if (_slugEdited) return;
    final slug = _nameCtrl.text
        .toLowerCase()
        .replaceAll(RegExp(r'[^a-z0-9\s\-]'), '')
        .trim()
        .replaceAll(RegExp(r'\s+'), '-');
    _slugCtrl.text = slug;
  }

  void _addPlan(String plan) {
    final trimmed = plan.trim().toLowerCase().replaceAll(' ', '-');
    if (trimmed.isEmpty || _plans.contains(trimmed)) return;
    setState(() {
      _plans.add(trimmed);
      _priceBase[trimmed] = TextEditingController(text: '0');
      _pricePerUser[trimmed] = TextEditingController(text: '0');
    });
    _planInputCtrl.clear();
  }

  void _removePlan(String plan) {
    setState(() {
      _plans.remove(plan);
      _priceBase.remove(plan)?.dispose();
      _pricePerUser.remove(plan)?.dispose();
    });
  }

  void _addModule() {
    setState(() {
      _modKeys.add(TextEditingController());
      _modNames.add(TextEditingController());
      _modDescs.add(TextEditingController());
    });
  }

  void _removeModule(int i) {
    setState(() {
      _modKeys.removeAt(i).dispose();
      _modNames.removeAt(i).dispose();
      _modDescs.removeAt(i).dispose();
    });
  }

  void _addApp() {
    setState(() {
      _appKeys.add(TextEditingController());
      _appNames.add(TextEditingController());
      _appDescs.add(TextEditingController());
    });
  }

  void _removeApp(int i) {
    setState(() {
      _appKeys.removeAt(i).dispose();
      _appNames.removeAt(i).dispose();
      _appDescs.removeAt(i).dispose();
    });
  }

  Future<void> _save() async {
    if (!_formKey.currentState!.validate()) return;
    setState(() => _isSaving = true);

    final modules = <Map<String, String>>[];
    for (int i = 0; i < _modKeys.length; i++) {
      modules.add({
        'key': _modKeys[i].text.trim(),
        'name': _modNames[i].text.trim(),
        'description': _modDescs[i].text.trim(),
      });
    }

    final apps = <Map<String, String>>[];
    for (int i = 0; i < _appKeys.length; i++) {
      apps.add({
        'key': _appKeys[i].text.trim(),
        'name': _appNames[i].text.trim(),
        'description': _appDescs[i].text.trim(),
      });
    }

    final pricingMap = <String, Map<String, dynamic>>{};
    for (final plan in _plans) {
      pricingMap[plan] = {
        'base_price': double.tryParse(_priceBase[plan]?.text ?? '0') ?? 0.0,
        'per_user_price':
            double.tryParse(_pricePerUser[plan]?.text ?? '0') ?? 0.0,
        'currency': 'IDR',
      };
    }

    final data = <String, dynamic>{
      'name': _nameCtrl.text.trim(),
      'slug': _slugCtrl.text.trim(),
      'description': _descCtrl.text.trim(),
      'is_active': _isActive,
      'available_plans': _plans,
      'base_pricing': pricingMap,
      'available_modules': modules,
      'available_apps': apps,
    };

    final error = await context.read<ProductCubit>().saveProduct(
          existing: widget.initialProduct,
          data: data,
        );

    if (!mounted) return;
    setState(() => _isSaving = false);

    if (error != null) {
      ScaffoldMessenger.of(context).showSnackBar(SnackBar(
        content: Text(error),
        backgroundColor: AppColors.errorBase,
        behavior: SnackBarBehavior.floating,
      ));
    } else {
      Navigator.pop(context, true);
    }
  }

  @override
  void dispose() {
    _nameCtrl.dispose();
    _slugCtrl.dispose();
    _descCtrl.dispose();
    _planInputCtrl.dispose();
    for (final ctrl in _priceBase.values) {
      ctrl.dispose();
    }
    for (final ctrl in _pricePerUser.values) {
      ctrl.dispose();
    }
    for (int i = 0; i < _modKeys.length; i++) {
      _modKeys[i].dispose();
      _modNames[i].dispose();
      _modDescs[i].dispose();
    }
    for (int i = 0; i < _appKeys.length; i++) {
      _appKeys[i].dispose();
      _appNames[i].dispose();
      _appDescs[i].dispose();
    }
    super.dispose();
  }

  // ── Build ──────────────────────────────────────────────────────────────

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: Text(_isEdit ? 'Edit Produk' : 'Tambah Produk'),
        backgroundColor: AppColors.primary700,
        foregroundColor: AppColors.textOnPrimary,
        elevation: 0,
        actions: [
          if (_isSaving)
            const Padding(
              padding: EdgeInsets.symmetric(horizontal: 16),
              child: Center(
                child: SizedBox(
                  width: 20,
                  height: 20,
                  child: CircularProgressIndicator(
                    strokeWidth: 2,
                    color: Colors.white,
                  ),
                ),
              ),
            )
          else
            TextButton.icon(
              onPressed: _save,
              icon: const Icon(Icons.check_rounded, color: Colors.white),
              label: const Text(
                'Simpan',
                style: TextStyle(
                  color: Colors.white,
                  fontWeight: FontWeight.w600,
                ),
              ),
            ),
        ],
      ),
      body: Form(
        key: _formKey,
        child: ListView(
          padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 16),
          children: [
            _buildBasicInfo(),
            const SizedBox(height: 16),
            _buildPlans(),
            if (_plans.isNotEmpty) ...[
              const SizedBox(height: 16),
              _buildPricing(),
            ],
            const SizedBox(height: 16),
            _buildModules(),
            const SizedBox(height: 16),
            _buildApps(),
            const SizedBox(height: 80),
          ],
        ),
      ),
    );
  }

  // ── Section builders ────────────────────────────────────────────────────

  Widget _buildBasicInfo() {
    return _sectionCard(
      title: 'Informasi Dasar',
      children: [
        _textField(
          ctrl: _nameCtrl,
          label: 'Nama Produk',
          hint: 'FlashERP',
          isRequired: true,
        ),
        const SizedBox(height: 12),
        _textField(
          ctrl: _slugCtrl,
          label: 'Slug',
          hint: 'flasherp',
          isRequired: true,
          helperText: 'Huruf kecil, angka, tanda hubung',
          onTap: () => _slugEdited = true,
          validator: (v) {
            if (v == null || v.isEmpty) return 'Slug wajib diisi';
            if (!RegExp(r'^[a-z0-9\-]+$').hasMatch(v)) {
              return 'Hanya huruf kecil, angka, dan tanda hubung';
            }
            return null;
          },
        ),
        const SizedBox(height: 12),
        _textField(
          ctrl: _descCtrl,
          label: 'Deskripsi',
          hint: 'Deskripsi singkat produk...',
          maxLines: 3,
        ),
        const SizedBox(height: 12),
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            const Text(
              'Status Aktif',
              style: TextStyle(
                fontSize: 14,
                fontWeight: FontWeight.w500,
                color: AppColors.textPrimary,
              ),
            ),
            Switch(
              value: _isActive,
              onChanged: (v) => setState(() => _isActive = v),
              activeTrackColor: AppColors.primary700,
            ),
          ],
        ),
      ],
    );
  }

  Widget _buildPlans() {
    final chips = <Widget>[];
    for (final plan in _plans) {
      chips.add(_planChip(plan));
    }
    return _sectionCard(
      title: 'Plan Tersedia',
      children: [
        Row(
          children: [
            Expanded(
              child: TextField(
                controller: _planInputCtrl,
                decoration: _inputDeco('Nama plan (mis: saas)', isDense: true),
                onSubmitted: _addPlan,
              ),
            ),
            const SizedBox(width: 8),
            ElevatedButton(
              onPressed: () => _addPlan(_planInputCtrl.text),
              style: ElevatedButton.styleFrom(
                backgroundColor: AppColors.primary700,
                foregroundColor: Colors.white,
                padding:
                    const EdgeInsets.symmetric(horizontal: 16, vertical: 10),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(10),
                ),
              ),
              child: const Text('Tambah'),
            ),
          ],
        ),
        const SizedBox(height: 8),
        if (_plans.isEmpty)
          const Padding(
            padding: EdgeInsets.symmetric(vertical: 4),
            child: Text(
              'Belum ada plan. Contoh: saas, on-premise',
              style:
                  TextStyle(color: AppColors.textSecondary, fontSize: 13),
            ),
          ),
        if (chips.isNotEmpty)
          Wrap(spacing: 8, runSpacing: 6, children: chips),
      ],
    );
  }

  Widget _planChip(String plan) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 5),
      decoration: BoxDecoration(
        color: AppColors.accent50,
        borderRadius: BorderRadius.circular(8),
        border: Border.all(color: AppColors.accent400),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Text(
            plan,
            style: const TextStyle(
              fontSize: 13,
              fontWeight: FontWeight.w500,
              color: AppColors.accent400,
            ),
          ),
          const SizedBox(width: 4),
          GestureDetector(
            onTap: () => _removePlan(plan),
            child: const Icon(Icons.close_rounded,
                size: 14, color: AppColors.accent400),
          ),
        ],
      ),
    );
  }

  Widget _buildPricing() {
    final cards = <Widget>[];
    for (final plan in _plans) {
      cards.add(_pricingCard(plan));
    }
    return _sectionCard(title: 'Harga per Plan', children: cards);
  }

  Widget _pricingCard(String plan) {
    return Container(
      margin: const EdgeInsets.only(bottom: 12),
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: AppColors.background,
        borderRadius: BorderRadius.circular(10),
        border: Border.all(color: AppColors.neutral200),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            plan,
            style: const TextStyle(
              fontSize: 13,
              fontWeight: FontWeight.w600,
              color: AppColors.primary600,
            ),
          ),
          const SizedBox(height: 10),
          Row(
            children: [
              Expanded(
                child: TextFormField(
                  controller: _priceBase[plan],
                  keyboardType: TextInputType.number,
                  decoration: _inputDeco('Harga Dasar (IDR)',
                      isDense: true, radius: 8),
                ),
              ),
              const SizedBox(width: 8),
              Expanded(
                child: TextFormField(
                  controller: _pricePerUser[plan],
                  keyboardType: TextInputType.number,
                  decoration:
                      _inputDeco('Per User (IDR)', isDense: true, radius: 8),
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildModules() {
    final cards = <Widget>[];
    for (int i = 0; i < _modKeys.length; i++) {
      cards.add(_itemCard(
        index: i,
        label: 'Modul',
        keyCtrl: _modKeys[i],
        nameCtrl: _modNames[i],
        descCtrl: _modDescs[i],
        onRemove: () => _removeModule(i),
      ));
    }
    return _sectionCard(
      title: 'Modul',
      children: [
        ...cards,
        const SizedBox(height: 4),
        _addButton('Tambah Modul', _addModule),
      ],
    );
  }

  Widget _buildApps() {
    final cards = <Widget>[];
    for (int i = 0; i < _appKeys.length; i++) {
      cards.add(_itemCard(
        index: i,
        label: 'Aplikasi',
        keyCtrl: _appKeys[i],
        nameCtrl: _appNames[i],
        descCtrl: _appDescs[i],
        onRemove: () => _removeApp(i),
      ));
    }
    return _sectionCard(
      title: 'Aplikasi',
      children: [
        ...cards,
        const SizedBox(height: 4),
        _addButton('Tambah Aplikasi', _addApp),
      ],
    );
  }

  Widget _itemCard({
    required int index,
    required String label,
    required TextEditingController keyCtrl,
    required TextEditingController nameCtrl,
    required TextEditingController descCtrl,
    required VoidCallback onRemove,
  }) {
    return Container(
      margin: const EdgeInsets.only(bottom: 12),
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: AppColors.background,
        borderRadius: BorderRadius.circular(10),
        border: Border.all(color: AppColors.neutral200),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Text(
                '$label ${index + 1}',
                style: const TextStyle(
                  fontSize: 12,
                  fontWeight: FontWeight.w600,
                  color: AppColors.textSecondary,
                ),
              ),
              const Spacer(),
              GestureDetector(
                onTap: onRemove,
                child: const Icon(Icons.delete_outline_rounded,
                    size: 18, color: AppColors.errorBase),
              ),
            ],
          ),
          const SizedBox(height: 8),
          TextFormField(
            controller: keyCtrl,
            decoration: _inputDeco('Key', hint: 'finance', isDense: true, radius: 8),
            validator: (v) =>
                (v == null || v.trim().isEmpty) ? 'Key wajib diisi' : null,
          ),
          const SizedBox(height: 8),
          TextFormField(
            controller: nameCtrl,
            decoration: _inputDeco('Nama',
                hint: 'Finance & Accounting', isDense: true, radius: 8),
            validator: (v) =>
                (v == null || v.trim().isEmpty) ? 'Nama wajib diisi' : null,
          ),
          const SizedBox(height: 8),
          TextFormField(
            controller: descCtrl,
            decoration:
                _inputDeco('Deskripsi', hint: 'Opsional', isDense: true, radius: 8),
          ),
        ],
      ),
    );
  }

  Widget _addButton(String label, VoidCallback onTap) {
    return OutlinedButton.icon(
      onPressed: onTap,
      icon: const Icon(Icons.add_rounded, size: 18),
      label: Text(label),
      style: OutlinedButton.styleFrom(
        foregroundColor: AppColors.primary700,
        side: const BorderSide(color: AppColors.primary700),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(10),
        ),
      ),
    );
  }

  // ── Shared widget helpers ───────────────────────────────────────────────

  Widget _sectionCard({
    required String title,
    required List<Widget> children,
  }) {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppColors.surface,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: AppColors.neutral200),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            title,
            style: const TextStyle(
              fontSize: 13,
              fontWeight: FontWeight.w700,
              color: AppColors.primary700,
              letterSpacing: 0.3,
            ),
          ),
          const SizedBox(height: 12),
          ...children,
        ],
      ),
    );
  }

  Widget _textField({
    required TextEditingController ctrl,
    required String label,
    bool isRequired = false,
    int maxLines = 1,
    String? hint,
    String? helperText,
    VoidCallback? onTap,
    String? Function(String?)? validator,
  }) {
    return TextFormField(
      controller: ctrl,
      maxLines: maxLines,
      onTap: onTap,
      decoration: _inputDeco(label, hint: hint, helperText: helperText)
          .copyWith(labelText: label + (isRequired ? ' *' : '')),
      validator: validator ??
          (isRequired
              ? (v) =>
                  (v == null || v.trim().isEmpty) ? '$label wajib diisi' : null
              : null),
    );
  }

  InputDecoration _inputDeco(
    String labelOrHint, {
    String? hint,
    String? helperText,
    bool isDense = false,
    double radius = 10,
  }) {
    return InputDecoration(
      labelText: labelOrHint,
      hintText: hint,
      helperText: helperText,
      hintStyle: const TextStyle(fontSize: 13),
      isDense: isDense,
      border: OutlineInputBorder(
        borderRadius: BorderRadius.circular(radius),
        borderSide: const BorderSide(color: AppColors.neutral300),
      ),
      enabledBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(radius),
        borderSide: const BorderSide(color: AppColors.neutral300),
      ),
      focusedBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(radius),
        borderSide: const BorderSide(color: AppColors.primary600, width: 1.5),
      ),
      errorBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(radius),
        borderSide: const BorderSide(color: AppColors.errorBase),
      ),
      focusedErrorBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(radius),
        borderSide: const BorderSide(color: AppColors.errorBase, width: 1.5),
      ),
      filled: true,
      fillColor: AppColors.surface,
      contentPadding: EdgeInsets.symmetric(
        horizontal: 14,
        vertical: isDense ? 10 : 12,
      ),
    );
  }
}
