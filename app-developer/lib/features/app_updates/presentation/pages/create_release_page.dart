import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import '../../../../core/constants/app_colors.dart';
import '../cubit/app_update_cubit.dart';

const _appIds = [
  'app-pos',
  'app-employee',
  'app-opname',
  'app-courier',
  'app-picking',
  'app-management',
  'app-customer',
  'app-supplier',
  'app-sales-person',
  'web-ui',
];

const _appLabels = <String, String>{
  'app-pos': 'Point of Sale',
  'app-employee': 'Karyawan',
  'app-opname': 'Opname Stok',
  'app-courier': 'Kurir',
  'app-picking': 'Picking',
  'app-management': 'Management',
  'app-customer': 'Customer',
  'app-supplier': 'Supplier',
  'app-sales-person': 'Sales Person',
  'web-ui': 'Web Admin',
};

class CreateReleasePage extends StatefulWidget {
  const CreateReleasePage({super.key});

  @override
  State<CreateReleasePage> createState() => _CreateReleasePageState();
}

class _CreateReleasePageState extends State<CreateReleasePage> {
  final _formKey = GlobalKey<FormState>();
  final _versionCtrl = TextEditingController();
  final _versionCodeCtrl = TextEditingController();
  final _downloadUrlCtrl = TextEditingController();
  final _releaseNotesCtrl = TextEditingController();

  String _selectedAppId = _appIds.first;
  bool _isMandatory = false;
  bool _submitting = false;

  @override
  void dispose() {
    _versionCtrl.dispose();
    _versionCodeCtrl.dispose();
    _downloadUrlCtrl.dispose();
    _releaseNotesCtrl.dispose();
    super.dispose();
  }

  Future<void> _submit() async {
    if (!_formKey.currentState!.validate()) return;

    setState(() => _submitting = true);

    await context.read<AppUpdateCubit>().createRelease(
          appId: _selectedAppId,
          version: _versionCtrl.text.trim(),
          versionCode: int.parse(_versionCodeCtrl.text.trim()),
          downloadUrl: _downloadUrlCtrl.text.trim(),
          releaseNotes: _releaseNotesCtrl.text.trim(),
          isMandatory: _isMandatory,
        );

    if (!mounted) return;
    setState(() => _submitting = false);

    final state = context.read<AppUpdateCubit>().state;
    if (state is ReleaseCreated) {
      Navigator.pop(context);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        backgroundColor: AppColors.primary800,
        foregroundColor: AppColors.white,
        title: const Text('Rilis Versi Baru'),
        leading: IconButton(
          icon: const Icon(Icons.close_rounded),
          onPressed: () => Navigator.pop(context),
        ),
      ),
      body: BlocListener<AppUpdateCubit, AppUpdateState>(
        listener: (context, state) {
          if (state is AppUpdateError) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                content: Text(state.message),
                backgroundColor: AppColors.error,
              ),
            );
          }
        },
        child: Form(
          key: _formKey,
          child: ListView(
            padding: const EdgeInsets.all(20),
            children: [
              // App selector
              _SectionLabel('Aplikasi'),
              Card(
                elevation: 0,
                color: AppColors.surface,
                shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                child: Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 16),
                  child: DropdownButtonFormField<String>(
                    value: _selectedAppId,
                    decoration: const InputDecoration(
                      labelText: 'Pilih Aplikasi',
                      border: InputBorder.none,
                    ),
                    items: _appIds
                        .map((id) => DropdownMenuItem(
                              value: id,
                              child: Text('${_appLabels[id] ?? id}  ($id)'),
                            ))
                        .toList(),
                    onChanged: (v) => setState(() => _selectedAppId = v!),
                  ),
                ),
              ),
              const SizedBox(height: 20),

              // Versi
              _SectionLabel('Versi'),
              Card(
                elevation: 0,
                color: AppColors.surface,
                shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                child: Padding(
                  padding: const EdgeInsets.all(16),
                  child: Column(
                    children: [
                      _InputField(
                        controller: _versionCtrl,
                        label: 'Versi (SemVer)',
                        hint: 'contoh: 1.2.3',
                        keyboardType: TextInputType.text,
                        validator: (v) {
                          if (v == null || v.isEmpty) return 'Wajib diisi';
                          final parts = v.split('.');
                          if (parts.length != 3) return 'Format: MAJOR.MINOR.PATCH (e.g. 1.2.3)';
                          return null;
                        },
                      ),
                      const Divider(height: 1),
                      _InputField(
                        controller: _versionCodeCtrl,
                        label: 'Version Code (integer)',
                        hint: 'contoh: 10203',
                        keyboardType: TextInputType.number,
                        inputFormatters: [FilteringTextInputFormatter.digitsOnly],
                        validator: (v) {
                          if (v == null || v.isEmpty) return 'Wajib diisi';
                          final code = int.tryParse(v);
                          if (code == null || code <= 0) return 'Harus bilangan positif';
                          return null;
                        },
                        helperText: 'Tip: MAJOR*10000 + MINOR*100 + PATCH (e.g. v1.2.3 = 10203)',
                      ),
                    ],
                  ),
                ),
              ),
              const SizedBox(height: 20),

              // Download URL
              _SectionLabel('URL Download'),
              Card(
                elevation: 0,
                color: AppColors.surface,
                shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                child: Padding(
                  padding: const EdgeInsets.all(16),
                  child: _InputField(
                    controller: _downloadUrlCtrl,
                    label: 'URL APK / Bundle',
                    hint: 'https://storage.example.com/releases/app-pos-1.2.3.apk',
                    keyboardType: TextInputType.url,
                    validator: (v) {
                      if (v == null || v.isEmpty) return 'Wajib diisi';
                      if (!v.startsWith('http')) return 'URL harus diawali http atau https';
                      return null;
                    },
                  ),
                ),
              ),
              const SizedBox(height: 20),

              // Catatan rilis
              _SectionLabel('Catatan Rilis'),
              Card(
                elevation: 0,
                color: AppColors.surface,
                shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                child: Padding(
                  padding: const EdgeInsets.all(16),
                  child: TextFormField(
                    controller: _releaseNotesCtrl,
                    maxLines: 5,
                    decoration: const InputDecoration(
                      labelText: 'Catatan perubahan',
                      hintText: '- Perbaikan bug login\n- Tambah fitur laporan\n- Performa lebih cepat',
                      border: InputBorder.none,
                      alignLabelWithHint: true,
                    ),
                  ),
                ),
              ),
              const SizedBox(height: 20),

              // Mandatory toggle
              Card(
                elevation: 0,
                color: AppColors.surface,
                shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                child: SwitchListTile(
                  value: _isMandatory,
                  onChanged: (v) => setState(() => _isMandatory = v),
                  activeColor: AppColors.primary700,
                  title: const Text(
                    'Update Wajib',
                    style: TextStyle(fontWeight: FontWeight.w600),
                  ),
                  subtitle: const Text(
                    'Aplikasi klien tidak bisa digunakan sebelum update',
                    style: TextStyle(fontSize: 12),
                  ),
                ),
              ),
              const SizedBox(height: 32),

              // Submit button
              ElevatedButton(
                style: ElevatedButton.styleFrom(
                  backgroundColor: AppColors.primary700,
                  foregroundColor: AppColors.white,
                  minimumSize: const Size.fromHeight(50),
                  shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                ),
                onPressed: _submitting ? null : _submit,
                child: _submitting
                    ? const SizedBox(
                        width: 20,
                        height: 20,
                        child: CircularProgressIndicator(
                          color: AppColors.white,
                          strokeWidth: 2,
                        ),
                      )
                    : const Text(
                        'Publikasikan Rilis',
                        style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold),
                      ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class _SectionLabel extends StatelessWidget {
  final String text;
  const _SectionLabel(this.text);

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 8),
      child: Text(
        text.toUpperCase(),
        style: const TextStyle(
          fontSize: 11,
          fontWeight: FontWeight.w700,
          color: AppColors.neutral500,
          letterSpacing: 0.8,
        ),
      ),
    );
  }
}

class _InputField extends StatelessWidget {
  final TextEditingController controller;
  final String label;
  final String hint;
  final TextInputType keyboardType;
  final String? Function(String?)? validator;
  final List<TextInputFormatter>? inputFormatters;
  final String? helperText;

  const _InputField({
    required this.controller,
    required this.label,
    required this.hint,
    this.keyboardType = TextInputType.text,
    this.validator,
    this.inputFormatters,
    this.helperText,
  });

  @override
  Widget build(BuildContext context) {
    return TextFormField(
      controller: controller,
      keyboardType: keyboardType,
      inputFormatters: inputFormatters,
      validator: validator,
      decoration: InputDecoration(
        labelText: label,
        hintText: hint,
        helperText: helperText,
        helperMaxLines: 2,
        border: InputBorder.none,
      ),
    );
  }
}
