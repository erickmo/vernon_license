import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:intl/intl.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/notifiers/pending_count_notifier.dart';
import '../../../../injection_container.dart';
import '../../../registrations/domain/entities/registration_entity.dart';
import '../../../registrations/presentation/cubit/registration_cubit.dart';
import '../../../registrations/presentation/pages/create_client_page.dart';
import '../../domain/entities/company_entity.dart';
import '../cubit/company_cubit.dart';

class ClientsPage extends StatelessWidget {
  const ClientsPage({super.key});

  @override
  Widget build(BuildContext context) {
    return MultiBlocProvider(
      providers: [
        BlocProvider(create: (_) => sl<CompanyCubit>()..load()),
        BlocProvider(
            create: (_) =>
                sl<RegistrationCubit>()..loadRegistrations(status: 'pending')),
      ],
      child: const _ClientsView(),
    );
  }
}

class _ClientsView extends StatefulWidget {
  const _ClientsView();

  @override
  State<_ClientsView> createState() => _ClientsViewState();
}

class _ClientsViewState extends State<_ClientsView>
    with SingleTickerProviderStateMixin {
  late final TabController _tabCtrl;
  final _searchCtrl = TextEditingController();
  bool _searchOpen = false;

  @override
  void initState() {
    super.initState();
    _tabCtrl = TabController(length: 2, vsync: this);
    _searchCtrl.addListener(() {
      context.read<CompanyCubit>().search(_searchCtrl.text);
    });
  }

  @override
  void dispose() {
    _tabCtrl.dispose();
    _searchCtrl.dispose();
    super.dispose();
  }

  void _openCreateClient() {
    final regCubit = context.read<RegistrationCubit>();
    Navigator.of(context).push(
      MaterialPageRoute(
        fullscreenDialog: true,
        builder: (_) => BlocProvider.value(
          value: regCubit,
          child: const CreateClientPage(),
        ),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return BlocListener<RegistrationCubit, RegistrationState>(
      listener: (context, state) {
        // Update pending badge
        if (state is RegistrationLoaded && state.activeFilter == 'pending') {
          sl<PendingCountNotifier>().update(state.items.length);
        }
        if (state is RegistrationEmpty && state.activeFilter == 'pending') {
          sl<PendingCountNotifier>().update(0);
        }
        // Snackbars
        if (state is RegistrationActionSuccess || state is ClientCreateSuccess) {
          final msg = state is RegistrationActionSuccess
              ? state.message
              : (state as ClientCreateSuccess).message;
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(content: Text(msg), backgroundColor: AppColors.success),
          );
          // Refresh klien aktif setelah create/approve
          context.read<CompanyCubit>().refresh();
        }
        if (state is RegistrationError) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
                content: Text(state.message),
                backgroundColor: AppColors.error),
          );
        }
        if (state is ClientCreateError) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
                content: Text(state.message),
                backgroundColor: AppColors.error),
          );
        }
      },
      child: Scaffold(
        backgroundColor: AppColors.background,
        appBar: AppBar(
          title: _searchOpen
              ? TextField(
                  controller: _searchCtrl,
                  autofocus: true,
                  decoration: const InputDecoration(
                    hintText: 'Cari nama / kode client...',
                    border: InputBorder.none,
                    hintStyle: TextStyle(color: AppColors.neutral400),
                  ),
                  style: const TextStyle(color: AppColors.neutral900, fontSize: 16),
                )
              : const Text('Client'),
          actions: [
            IconButton(
              icon: Icon(_searchOpen ? Icons.close_rounded : Icons.search_rounded),
              onPressed: () {
                setState(() {
                  _searchOpen = !_searchOpen;
                  if (!_searchOpen) {
                    _searchCtrl.clear();
                    context.read<CompanyCubit>().refresh();
                  }
                });
              },
              tooltip: _searchOpen ? 'Tutup pencarian' : 'Cari client',
            ),
            IconButton(
              icon: const Icon(Icons.refresh_rounded),
              onPressed: () {
                context.read<CompanyCubit>().refresh();
                context
                    .read<RegistrationCubit>()
                    .loadRegistrations(status: 'pending');
              },
              tooltip: 'Muat ulang',
            ),
          ],
          bottom: TabBar(
            controller: _tabCtrl,
            labelColor: AppColors.primary700,
            unselectedLabelColor: AppColors.neutral500,
            indicatorColor: AppColors.primary700,
            indicatorWeight: 2.5,
            tabs: [
              const Tab(text: 'Klien Aktif'),
              Tab(
                child: ValueListenableBuilder<int>(
                  valueListenable: sl<PendingCountNotifier>(),
                  builder: (_, count, __) => Row(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      const Text('Pengajuan'),
                      if (count > 0) ...[
                        const SizedBox(width: 6),
                        Container(
                          padding: const EdgeInsets.symmetric(
                              horizontal: 6, vertical: 1),
                          decoration: BoxDecoration(
                            color: AppColors.error,
                            borderRadius: BorderRadius.circular(10),
                          ),
                          child: Text(
                            '$count',
                            style: const TextStyle(
                                color: AppColors.white,
                                fontSize: 11,
                                fontWeight: FontWeight.w700),
                          ),
                        ),
                      ],
                    ],
                  ),
                ),
              ),
            ],
          ),
        ),
        body: TabBarView(
          controller: _tabCtrl,
          children: const [
            _ActiveClientsTab(),
            _RegistrationsTab(),
          ],
        ),
        floatingActionButton: FloatingActionButton.extended(
          onPressed: _openCreateClient,
          backgroundColor: AppColors.primary700,
          foregroundColor: AppColors.white,
          icon: const Icon(Icons.add_business_rounded),
          label: const Text('Client Baru'),
        ),
      ),
    );
  }
}

// ─────────────────────────────────────────────────────────────────────────────
// Tab 1: Klien Aktif
// ─────────────────────────────────────────────────────────────────────────────

class _ActiveClientsTab extends StatelessWidget {
  const _ActiveClientsTab();

  @override
  Widget build(BuildContext context) {
    return BlocBuilder<CompanyCubit, CompanyState>(
      builder: (context, state) {
        if (state is CompanyLoading || state is CompanyInitial) {
          return const Center(child: CircularProgressIndicator());
        }
        if (state is CompanyError) {
          return _ErrorView(
            message: state.message,
            onRetry: () => context.read<CompanyCubit>().refresh(),
          );
        }
        if (state is CompanyEmpty) {
          return _EmptyClients(search: state.search);
        }
        if (state is CompanyLoaded) {
          return _CompanyList(items: state.items);
        }
        return const SizedBox.shrink();
      },
    );
  }
}

class _CompanyList extends StatelessWidget {
  final List<CompanyEntity> items;

  const _CompanyList({required this.items});

  @override
  Widget build(BuildContext context) {
    return ListView.separated(
      padding: const EdgeInsets.fromLTRB(16, 16, 16, 100),
      itemCount: items.length,
      separatorBuilder: (_, __) => const SizedBox(height: 10),
      itemBuilder: (_, i) => _CompanyCard(company: items[i]),
    );
  }
}

class _CompanyCard extends StatelessWidget {
  final CompanyEntity company;

  const _CompanyCard({required this.company});

  @override
  Widget build(BuildContext context) {
    return Card(
      margin: EdgeInsets.zero,
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // ── Header ───────────────────────────────────────────────────
            Row(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                // Code badge
                Container(
                  padding:
                      const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                  decoration: BoxDecoration(
                    color: AppColors.primary700,
                    borderRadius: BorderRadius.circular(6),
                  ),
                  child: Text(
                    company.code,
                    style: const TextStyle(
                      color: AppColors.white,
                      fontSize: 11,
                      fontWeight: FontWeight.w700,
                      fontFamily: 'monospace',
                    ),
                  ),
                ),
                const SizedBox(width: 10),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        company.name,
                        style: const TextStyle(
                          fontSize: 15,
                          fontWeight: FontWeight.w700,
                          color: AppColors.neutral900,
                        ),
                      ),
                      if (company.companyType != null)
                        Text(
                          company.companyType!,
                          style: const TextStyle(
                              fontSize: 12, color: AppColors.neutral500),
                        ),
                    ],
                  ),
                ),
                // Status badge
                _StatusDot(isActive: company.isActive),
              ],
            ),

            const SizedBox(height: 12),
            const Divider(height: 1),
            const SizedBox(height: 12),

            // ── Stats row ─────────────────────────────────────────────────
            Row(
              children: [
                _StatChip(
                  icon: Icons.extension_rounded,
                  label: '${company.modules.length} modul',
                  color: AppColors.primary700,
                  bg: AppColors.primary50,
                ),
                const SizedBox(width: 8),
                _StatChip(
                  icon: Icons.apps_rounded,
                  label: '${company.apps.length} aplikasi',
                  color: AppColors.accent400,
                  bg: AppColors.accent50,
                ),
                const Spacer(),
                Text(
                  DateFormat('dd MMM yyyy', 'id_ID')
                      .format(company.createdAt.toLocal()),
                  style: const TextStyle(
                      fontSize: 11, color: AppColors.neutral400),
                ),
              ],
            ),

            // ── Kontak snippet ────────────────────────────────────────────
            if (company.email != null || company.phone != null) ...[
              const SizedBox(height: 10),
              Wrap(
                spacing: 12,
                runSpacing: 4,
                children: [
                  if (company.email != null)
                    _ContactChip(
                        icon: Icons.email_outlined, text: company.email!),
                  if (company.phone != null)
                    _ContactChip(
                        icon: Icons.phone_outlined, text: company.phone!),
                ],
              ),
            ],

            // ── Module chips ──────────────────────────────────────────────
            if (company.modules.isNotEmpty) ...[
              const SizedBox(height: 10),
              Wrap(
                spacing: 6,
                runSpacing: 6,
                children: company.modules
                    .map((m) => _MiniChip(label: _moduleName(m)))
                    .toList(),
              ),
            ],
          ],
        ),
      ),
    );
  }

  static String _moduleName(String id) {
    return const {
      'accounting': 'Akuntansi',
      'sales': 'Penjualan',
      'purchasing': 'Pembelian',
      'stock': 'Inventori',
      'courier': 'Kurir',
      'pos': 'POS',
      'hrm': 'HR & Payroll',
      'fixed_assets': 'Aset Tetap',
      'budgeting': 'Anggaran',
      'project': 'Proyek',
    }[id] ??
        id;
  }
}

class _StatusDot extends StatelessWidget {
  final bool isActive;

  const _StatusDot({required this.isActive});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      decoration: BoxDecoration(
        color: isActive
            ? AppColors.success.withValues(alpha: 0.12)
            : AppColors.neutral200,
        borderRadius: BorderRadius.circular(20),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Container(
            width: 6,
            height: 6,
            decoration: BoxDecoration(
              color: isActive ? AppColors.success : AppColors.neutral400,
              shape: BoxShape.circle,
            ),
          ),
          const SizedBox(width: 4),
          Text(
            isActive ? 'Aktif' : 'Nonaktif',
            style: TextStyle(
              fontSize: 11,
              fontWeight: FontWeight.w600,
              color: isActive ? AppColors.successBase : AppColors.neutral500,
            ),
          ),
        ],
      ),
    );
  }
}

class _StatChip extends StatelessWidget {
  final IconData icon;
  final String label;
  final Color color;
  final Color bg;

  const _StatChip(
      {required this.icon,
      required this.label,
      required this.color,
      required this.bg});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 5),
      decoration: BoxDecoration(
        color: bg,
        borderRadius: BorderRadius.circular(20),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(icon, size: 13, color: color),
          const SizedBox(width: 4),
          Text(label,
              style: TextStyle(
                  fontSize: 12, fontWeight: FontWeight.w500, color: color)),
        ],
      ),
    );
  }
}

class _ContactChip extends StatelessWidget {
  final IconData icon;
  final String text;

  const _ContactChip({required this.icon, required this.text});

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        Icon(icon, size: 13, color: AppColors.neutral400),
        const SizedBox(width: 4),
        Text(text,
            style:
                const TextStyle(fontSize: 12, color: AppColors.neutral600)),
      ],
    );
  }
}

class _MiniChip extends StatelessWidget {
  final String label;

  const _MiniChip({required this.label});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 3),
      decoration: BoxDecoration(
        color: AppColors.neutral100,
        borderRadius: BorderRadius.circular(20),
        border: Border.all(color: AppColors.neutral200),
      ),
      child: Text(label,
          style: const TextStyle(
              fontSize: 11,
              color: AppColors.neutral600,
              fontWeight: FontWeight.w500)),
    );
  }
}

class _EmptyClients extends StatelessWidget {
  final String search;

  const _EmptyClients({required this.search});

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          const Icon(Icons.business_outlined,
              size: 64, color: AppColors.neutral300),
          const SizedBox(height: 12),
          Text(
            search.isEmpty
                ? 'Belum ada client yang terdaftar'
                : 'Tidak ada client dengan kata kunci "$search"',
            textAlign: TextAlign.center,
            style:
                const TextStyle(fontSize: 15, color: AppColors.neutral500),
          ),
          if (search.isEmpty) ...[
            const SizedBox(height: 8),
            const Text(
              'Buat client baru atau setujui pengajuan registrasi',
              textAlign: TextAlign.center,
              style: TextStyle(fontSize: 13, color: AppColors.neutral400),
            ),
          ],
          const SizedBox(height: 20),
          TextButton.icon(
            onPressed: () => context.read<CompanyCubit>().refresh(),
            icon: const Icon(Icons.refresh_rounded),
            label: const Text('Muat Ulang'),
          ),
        ],
      ),
    );
  }
}

// ─────────────────────────────────────────────────────────────────────────────
// Tab 2: Pengajuan Registrasi
// ─────────────────────────────────────────────────────────────────────────────

class _RegistrationsTab extends StatelessWidget {
  const _RegistrationsTab();

  @override
  Widget build(BuildContext context) {
    return BlocBuilder<RegistrationCubit, RegistrationState>(
      builder: (context, state) {
        final activeFilter = _getFilter(state);

        return Column(
          children: [
            _RegFilterBar(activeFilter: activeFilter),
            Expanded(child: _buildBody(context, state)),
          ],
        );
      },
    );
  }

  String? _getFilter(RegistrationState state) {
    if (state is RegistrationLoaded) return state.activeFilter;
    if (state is RegistrationEmpty) return state.activeFilter;
    if (state is RegistrationActionLoading) return state.activeFilter;
    if (state is RegistrationActionSuccess) return state.activeFilter;
    if (state is ClientCreateLoading) return state.activeFilter;
    if (state is ClientCreateSuccess) return state.activeFilter;
    if (state is ClientCreateError) return state.activeFilter;
    return 'pending';
  }

  Widget _buildBody(BuildContext context, RegistrationState state) {
    if (state is RegistrationLoading || state is RegistrationInitial) {
      return const Center(child: CircularProgressIndicator());
    }
    final items = _getItems(state);
    if (items != null) {
      final loading = state is RegistrationActionLoading ||
          state is ClientCreateLoading;
      if (items.isEmpty) {
        return _RegEmpty(filter: _getFilter(state) ?? 'pending');
      }
      return _RegList(items: items, isLoading: loading);
    }
    if (state is RegistrationError) {
      return _ErrorView(
        message: state.message,
        onRetry: () => context.read<RegistrationCubit>().loadRegistrations(),
      );
    }
    return const SizedBox.shrink();
  }

  List<RegistrationEntity>? _getItems(RegistrationState state) {
    if (state is RegistrationLoaded) return state.items;
    if (state is RegistrationActionLoading) return state.items;
    if (state is RegistrationActionSuccess) return state.items;
    if (state is RegistrationEmpty) return const [];
    if (state is ClientCreateLoading) return state.items;
    if (state is ClientCreateSuccess) return state.items;
    if (state is ClientCreateError) return state.items;
    return null;
  }
}

class _RegFilterBar extends StatelessWidget {
  final String? activeFilter;

  const _RegFilterBar({this.activeFilter});

  @override
  Widget build(BuildContext context) {
    const filters = [
      ('pending', 'Menunggu'),
      ('approved', 'Disetujui'),
      ('rejected', 'Ditolak'),
      (null, 'Semua'),
    ];

    return Container(
      color: AppColors.surface,
      padding: const EdgeInsets.fromLTRB(16, 10, 16, 10),
      child: SingleChildScrollView(
        scrollDirection: Axis.horizontal,
        child: Row(
          children: filters.map((f) {
            final isActive = activeFilter == f.$1;
            return Padding(
              padding: const EdgeInsets.only(right: 8),
              child: FilterChip(
                label: Text(f.$2),
                selected: isActive,
                onSelected: (_) => context
                    .read<RegistrationCubit>()
                    .loadRegistrations(status: f.$1),
                selectedColor: AppColors.primary100,
                checkmarkColor: AppColors.primary700,
                labelStyle: TextStyle(
                  color:
                      isActive ? AppColors.primary700 : AppColors.neutral700,
                  fontWeight:
                      isActive ? FontWeight.w600 : FontWeight.normal,
                  fontSize: 13,
                ),
                side: BorderSide(
                  color: isActive
                      ? AppColors.primary700
                      : AppColors.neutral300,
                ),
              ),
            );
          }).toList(),
        ),
      ),
    );
  }
}

class _RegList extends StatelessWidget {
  final List<RegistrationEntity> items;
  final bool isLoading;

  const _RegList({required this.items, this.isLoading = false});

  @override
  Widget build(BuildContext context) {
    return Stack(
      children: [
        ListView.separated(
          padding: const EdgeInsets.fromLTRB(16, 16, 16, 100),
          itemCount: items.length,
          separatorBuilder: (_, __) => const SizedBox(height: 10),
          itemBuilder: (_, i) => _RegCard(item: items[i]),
        ),
        if (isLoading)
          Container(
            color: Colors.black.withValues(alpha: 0.05),
            child: const Center(child: CircularProgressIndicator()),
          ),
      ],
    );
  }
}

class _RegCard extends StatelessWidget {
  final RegistrationEntity item;

  const _RegCard({required this.item});

  @override
  Widget build(BuildContext context) {
    return Card(
      margin: EdgeInsets.zero,
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Expanded(
                  child: Text(
                    item.nama,
                    style: const TextStyle(
                        fontSize: 15,
                        fontWeight: FontWeight.w700,
                        color: AppColors.neutral900),
                  ),
                ),
                _RegStatusBadge(status: item.status),
              ],
            ),
            const SizedBox(height: 8),
            if (item.npwp != null)
              _RegInfoRow(
                  icon: Icons.receipt_outlined, text: 'NPWP: ${item.npwp}'),
            if (item.email != null)
              _RegInfoRow(icon: Icons.email_outlined, text: item.email!),
            if (item.telepon != null)
              _RegInfoRow(
                  icon: Icons.phone_outlined, text: item.telepon!),
            if (item.alamat != null)
              _RegInfoRow(
                  icon: Icons.location_on_outlined, text: item.alamat!),
            if (item.catatan != null && item.catatan!.isNotEmpty)
              _RegInfoRow(
                  icon: Icons.notes_outlined,
                  text: 'Catatan: ${item.catatan}'),
            const SizedBox(height: 6),
            Text(
              'Diajukan: ${DateFormat('dd MMM yyyy, HH:mm', 'id_ID').format(item.createdAt.toLocal())}',
              style: const TextStyle(
                  fontSize: 12, color: AppColors.neutral400),
            ),
            if (item.isPending) ...[
              const SizedBox(height: 14),
              const Divider(height: 1),
              const SizedBox(height: 12),
              Row(
                children: [
                  Expanded(
                    child: OutlinedButton.icon(
                      onPressed: () => _showRejectDialog(context),
                      icon: const Icon(Icons.close_rounded, size: 16),
                      label: const Text('Tolak'),
                      style: OutlinedButton.styleFrom(
                        foregroundColor: AppColors.error,
                        side: const BorderSide(color: AppColors.error),
                      ),
                    ),
                  ),
                  const SizedBox(width: 10),
                  Expanded(
                    child: ElevatedButton.icon(
                      onPressed: () => _showApproveDialog(context),
                      icon: const Icon(Icons.check_rounded, size: 16),
                      label: const Text('Setujui'),
                      style: ElevatedButton.styleFrom(
                        backgroundColor: AppColors.success,
                        minimumSize: const Size(0, 40),
                      ),
                    ),
                  ),
                ],
              ),
            ],
          ],
        ),
      ),
    );
  }

  void _showApproveDialog(BuildContext context) {
    showDialog(
      context: context,
      builder: (_) => _ApproveDialog(
        registration: item,
        onApprove: (code, name) => context.read<RegistrationCubit>().approve(
              id: item.id,
              companyCode: code,
              companyName: name,
            ),
      ),
    );
  }

  void _showRejectDialog(BuildContext context) {
    showDialog(
      context: context,
      builder: (_) => _RejectDialog(
        registration: item,
        onReject: (reason) => context.read<RegistrationCubit>().reject(
              id: item.id,
              reason: reason,
            ),
      ),
    );
  }
}

class _RegInfoRow extends StatelessWidget {
  final IconData icon;
  final String text;

  const _RegInfoRow({required this.icon, required this.text});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 4),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Icon(icon, size: 14, color: AppColors.neutral400),
          const SizedBox(width: 6),
          Expanded(
            child: Text(text,
                style:
                    const TextStyle(fontSize: 13, color: AppColors.neutral700)),
          ),
        ],
      ),
    );
  }
}

class _RegStatusBadge extends StatelessWidget {
  final String status;

  const _RegStatusBadge({required this.status});

  @override
  Widget build(BuildContext context) {
    final (label, color, bg) = switch (status) {
      'pending' => ('Menunggu', AppColors.warning, const Color(0xFFFFF7E0)),
      'approved' => ('Disetujui', AppColors.successBase, const Color(0xFFE8FFF0)),
      'rejected' => ('Ditolak', AppColors.error, const Color(0xFFFFEEEE)),
      _ => ('Unknown', AppColors.neutral500, AppColors.neutral100),
    };
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
      decoration:
          BoxDecoration(color: bg, borderRadius: BorderRadius.circular(20)),
      child: Text(label,
          style: TextStyle(
              color: color, fontSize: 12, fontWeight: FontWeight.w600)),
    );
  }
}

class _RegEmpty extends StatelessWidget {
  final String filter;

  const _RegEmpty({required this.filter});

  @override
  Widget build(BuildContext context) {
    final label = switch (filter) {
      'pending' => 'Tidak ada pengajuan yang menunggu',
      'approved' => 'Belum ada registrasi yang disetujui',
      'rejected' => 'Belum ada registrasi yang ditolak',
      _ => 'Belum ada data pengajuan',
    };
    return Center(
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          const Icon(Icons.inbox_rounded,
              size: 64, color: AppColors.neutral300),
          const SizedBox(height: 12),
          Text(label,
              style:
                  const TextStyle(fontSize: 15, color: AppColors.neutral500)),
          const SizedBox(height: 16),
          TextButton.icon(
            onPressed: () => context
                .read<RegistrationCubit>()
                .loadRegistrations(
                    status: filter == 'semua' ? null : filter),
            icon: const Icon(Icons.refresh_rounded),
            label: const Text('Muat Ulang'),
          ),
        ],
      ),
    );
  }
}

// ─────────────────────────────────────────────────────────────────────────────
// Approve / Reject Dialogs
// ─────────────────────────────────────────────────────────────────────────────

class _ApproveDialog extends StatefulWidget {
  final RegistrationEntity registration;
  final void Function(String code, String name) onApprove;

  const _ApproveDialog({required this.registration, required this.onApprove});

  @override
  State<_ApproveDialog> createState() => _ApproveDialogState();
}

class _ApproveDialogState extends State<_ApproveDialog> {
  final _formKey = GlobalKey<FormState>();
  final _codeCtrl = TextEditingController();
  final _nameCtrl = TextEditingController();

  @override
  void initState() {
    super.initState();
    _nameCtrl.text = widget.registration.nama;
  }

  @override
  void dispose() {
    _codeCtrl.dispose();
    _nameCtrl.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      title: const Text('Setujui Registrasi'),
      content: Form(
        key: _formKey,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'Anda akan menyetujui:\n${widget.registration.nama}',
              style:
                  const TextStyle(fontSize: 14, color: AppColors.neutral700),
            ),
            const SizedBox(height: 4),
            const Text(
              'Tindakan ini akan membuat perusahaan baru di sistem.',
              style: TextStyle(fontSize: 12, color: AppColors.neutral500),
            ),
            const SizedBox(height: 20),
            TextFormField(
              controller: _codeCtrl,
              textCapitalization: TextCapitalization.characters,
              decoration: const InputDecoration(
                labelText: 'Kode Perusahaan *',
                hintText: 'misal: CUST-001',
              ),
              validator: (v) => v == null || v.trim().isEmpty
                  ? 'Kode wajib diisi'
                  : null,
            ),
            const SizedBox(height: 12),
            TextFormField(
              controller: _nameCtrl,
              decoration: const InputDecoration(
                labelText: 'Nama Resmi Perusahaan *',
              ),
              validator: (v) => v == null || v.trim().isEmpty
                  ? 'Nama wajib diisi'
                  : null,
            ),
          ],
        ),
      ),
      actions: [
        TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Batal')),
        ElevatedButton(
          onPressed: () {
            if (!_formKey.currentState!.validate()) return;
            Navigator.pop(context);
            widget.onApprove(
                _codeCtrl.text.trim().toUpperCase(), _nameCtrl.text.trim());
          },
          style: ElevatedButton.styleFrom(backgroundColor: AppColors.success),
          child: const Text('Setujui'),
        ),
      ],
    );
  }
}

class _RejectDialog extends StatefulWidget {
  final RegistrationEntity registration;
  final void Function(String reason) onReject;

  const _RejectDialog({required this.registration, required this.onReject});

  @override
  State<_RejectDialog> createState() => _RejectDialogState();
}

class _RejectDialogState extends State<_RejectDialog> {
  final _formKey = GlobalKey<FormState>();
  final _reasonCtrl = TextEditingController();

  @override
  void dispose() {
    _reasonCtrl.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      title: const Text('Tolak Registrasi'),
      content: Form(
        key: _formKey,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Text('Anda akan menolak:\n${widget.registration.nama}',
                style: const TextStyle(
                    fontSize: 14, color: AppColors.neutral700)),
            const SizedBox(height: 16),
            TextFormField(
              controller: _reasonCtrl,
              maxLines: 3,
              decoration: const InputDecoration(
                labelText: 'Alasan Penolakan *',
                hintText: 'Jelaskan alasan penolakan...',
                alignLabelWithHint: true,
              ),
              validator: (v) =>
                  v == null || v.trim().isEmpty ? 'Alasan wajib diisi' : null,
            ),
          ],
        ),
      ),
      actions: [
        TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Batal')),
        ElevatedButton(
          onPressed: () {
            if (!_formKey.currentState!.validate()) return;
            Navigator.pop(context);
            widget.onReject(_reasonCtrl.text.trim());
          },
          style: ElevatedButton.styleFrom(backgroundColor: AppColors.error),
          child: const Text('Tolak'),
        ),
      ],
    );
  }
}

// ─────────────────────────────────────────────────────────────────────────────
// Shared
// ─────────────────────────────────────────────────────────────────────────────

class _ErrorView extends StatelessWidget {
  final String message;
  final VoidCallback onRetry;

  const _ErrorView({required this.message, required this.onRetry});

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          const Icon(Icons.error_outline, size: 48, color: AppColors.error),
          const SizedBox(height: 12),
          Text(message,
              textAlign: TextAlign.center,
              style: const TextStyle(color: AppColors.neutral700)),
          const SizedBox(height: 16),
          ElevatedButton.icon(
            onPressed: onRetry,
            icon: const Icon(Icons.refresh_rounded),
            label: const Text('Coba Lagi'),
          ),
        ],
      ),
    );
  }
}
