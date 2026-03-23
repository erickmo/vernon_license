import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../features/clients/domain/entities/company_entity.dart';
import '../../../../features/clients/domain/usecases/list_companies_usecase.dart';
import '../../../../injection_container.dart' show sl;
import '../../domain/entities/app_release_entity.dart';
import '../cubit/app_update_cubit.dart';

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

class PushUpdatePage extends StatefulWidget {
  final AppReleaseEntity release;

  const PushUpdatePage({super.key, required this.release});

  @override
  State<PushUpdatePage> createState() => _PushUpdatePageState();
}

class _PushUpdatePageState extends State<PushUpdatePage> {
  bool _loadingCompanies = true;
  List<CompanyEntity> _companies = [];
  String? _companiesError;

  CompanyEntity? _selectedCompany;
  bool _forceUpdate = false;
  bool _submitting = false;

  @override
  void initState() {
    super.initState();
    _loadCompanies();
  }

  Future<void> _loadCompanies() async {
    final useCase = sl<ListCompaniesUseCase>();
    final result = await useCase();
    if (!mounted) return;

    result.fold(
      (failure) => setState(() {
        _companiesError = failure.message;
        _loadingCompanies = false;
      }),
      (companies) => setState(() {
        _companies = companies.where((c) => c.isActive).toList();
        _loadingCompanies = false;
      }),
    );
  }

  Future<void> _submit() async {
    if (_selectedCompany == null) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('Pilih klien tujuan terlebih dahulu'),
          backgroundColor: AppColors.warning,
        ),
      );
      return;
    }

    setState(() => _submitting = true);

    await context.read<AppUpdateCubit>().pushUpdate(
          companyId: _selectedCompany!.id,
          appId: widget.release.appId,
          versionCode: widget.release.versionCode,
          forceUpdate: _forceUpdate,
        );

    if (!mounted) return;
    setState(() => _submitting = false);

    final state = context.read<AppUpdateCubit>().state;
    if (state is UpdatePushed) {
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
        title: const Text('Push Update ke Klien'),
        leading: IconButton(
          icon: const Icon(Icons.close_rounded),
          onPressed: () => Navigator.pop(context),
        ),
      ),
      body: BlocListener<AppUpdateCubit, AppUpdateState>(
        listener: (context, state) {
          if (state is PushUpdateError) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                content: Text(state.message),
                backgroundColor: AppColors.error,
              ),
            );
          }
        },
        child: ListView(
          padding: const EdgeInsets.all(20),
          children: [
            // Info rilis
            _ReleaseInfoCard(release: widget.release),
            const SizedBox(height: 20),

            // Pilih klien
            _SectionLabel('Klien Tujuan'),
            _CompanySelector(
              loading: _loadingCompanies,
              error: _companiesError,
              companies: _companies,
              selected: _selectedCompany,
              onChanged: (c) => setState(() => _selectedCompany = c),
            ),
            const SizedBox(height: 20),

            // Force update toggle
            Card(
              elevation: 0,
              color: AppColors.surface,
              shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
              child: SwitchListTile(
                value: _forceUpdate,
                onChanged: (v) => setState(() => _forceUpdate = v),
                activeColor: AppColors.errorBase,
                title: const Text(
                  'Update Wajib (Force)',
                  style: TextStyle(fontWeight: FontWeight.w600),
                ),
                subtitle: const Text(
                  'Aplikasi klien tidak bisa digunakan sebelum update selesai',
                  style: TextStyle(fontSize: 12),
                ),
              ),
            ),
            if (_forceUpdate) ...[
              const SizedBox(height: 8),
              Container(
                padding: const EdgeInsets.all(12),
                decoration: BoxDecoration(
                  color: AppColors.errorLight,
                  borderRadius: BorderRadius.circular(8),
                ),
                child: const Row(
                  children: [
                    Icon(Icons.warning_amber_rounded, color: AppColors.errorBase, size: 18),
                    SizedBox(width: 8),
                    Expanded(
                      child: Text(
                        'Force update akan memblokir akses klien ke aplikasi sampai update selesai diinstall.',
                        style: TextStyle(fontSize: 12, color: AppColors.errorBase),
                      ),
                    ),
                  ],
                ),
              ),
            ],
            const SizedBox(height: 32),

            ElevatedButton.icon(
              style: ElevatedButton.styleFrom(
                backgroundColor: AppColors.primary700,
                foregroundColor: AppColors.white,
                minimumSize: const Size.fromHeight(50),
                shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
              ),
              onPressed: _submitting ? null : _submit,
              icon: _submitting
                  ? const SizedBox(
                      width: 18,
                      height: 18,
                      child: CircularProgressIndicator(
                        color: AppColors.white,
                        strokeWidth: 2,
                      ),
                    )
                  : const Icon(Icons.send_rounded),
              label: Text(
                _submitting ? 'Memproses...' : 'Push Update',
                style: const TextStyle(fontSize: 16, fontWeight: FontWeight.bold),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _ReleaseInfoCard extends StatelessWidget {
  final AppReleaseEntity release;

  const _ReleaseInfoCard({required this.release});

  @override
  Widget build(BuildContext context) {
    return Card(
      elevation: 0,
      color: AppColors.primary50,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(12),
        side: const BorderSide(color: AppColors.primary100),
      ),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                const Icon(Icons.system_update_rounded, color: AppColors.primary700, size: 20),
                const SizedBox(width: 8),
                Text(
                  _appLabels[release.appId] ?? release.appId,
                  style: const TextStyle(
                    fontWeight: FontWeight.w600,
                    color: AppColors.primary700,
                  ),
                ),
                const Spacer(),
                if (release.isMandatory)
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 3),
                    decoration: BoxDecoration(
                      color: AppColors.errorLight,
                      borderRadius: BorderRadius.circular(6),
                    ),
                    child: const Text(
                      'WAJIB',
                      style: TextStyle(
                        fontSize: 10,
                        fontWeight: FontWeight.w700,
                        color: AppColors.errorBase,
                      ),
                    ),
                  ),
              ],
            ),
            const SizedBox(height: 8),
            Text(
              'v${release.version} (${release.versionCode})',
              style: const TextStyle(
                fontSize: 20,
                fontWeight: FontWeight.bold,
                color: AppColors.primary800,
                fontFamily: 'JetBrainsMono',
              ),
            ),
            if (release.releaseNotes.isNotEmpty) ...[
              const SizedBox(height: 6),
              Text(
                release.releaseNotes,
                style: const TextStyle(fontSize: 13, color: AppColors.primary600),
                maxLines: 2,
                overflow: TextOverflow.ellipsis,
              ),
            ],
          ],
        ),
      ),
    );
  }
}

class _CompanySelector extends StatelessWidget {
  final bool loading;
  final String? error;
  final List<CompanyEntity> companies;
  final CompanyEntity? selected;
  final ValueChanged<CompanyEntity?> onChanged;

  const _CompanySelector({
    required this.loading,
    this.error,
    required this.companies,
    required this.selected,
    required this.onChanged,
  });

  @override
  Widget build(BuildContext context) {
    if (loading) {
      return const Card(
        child: Padding(
          padding: EdgeInsets.all(24),
          child: Center(child: CircularProgressIndicator()),
        ),
      );
    }

    if (error != null) {
      return Card(
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Text(error!, style: const TextStyle(color: AppColors.error)),
        ),
      );
    }

    if (companies.isEmpty) {
      return Card(
        elevation: 0,
        color: AppColors.surface,
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
        child: const Padding(
          padding: EdgeInsets.all(16),
          child: Text(
            'Belum ada klien aktif',
            style: TextStyle(color: AppColors.neutral500),
          ),
        ),
      );
    }

    return Card(
      elevation: 0,
      color: AppColors.surface,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 16),
        child: DropdownButtonFormField<CompanyEntity>(
          value: selected,
          decoration: const InputDecoration(
            labelText: 'Pilih Klien',
            border: InputBorder.none,
          ),
          items: companies
              .map(
                (c) => DropdownMenuItem(
                  value: c,
                  child: Text(
                    '${c.name}  (${c.code})',
                    overflow: TextOverflow.ellipsis,
                  ),
                ),
              )
              .toList(),
          onChanged: onChanged,
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
