import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:intl/intl.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/notifiers/pending_count_notifier.dart';
import '../../../../injection_container.dart';
import '../../domain/entities/registration_entity.dart';
import '../cubit/registration_cubit.dart';
import 'create_client_page.dart';

class RegistrationsPage extends StatelessWidget {
  const RegistrationsPage({super.key});

  @override
  Widget build(BuildContext context) {
    return BlocProvider(
      create: (_) => sl<RegistrationCubit>()..loadRegistrations(status: 'pending'),
      child: const _RegistrationsView(),
    );
  }
}

class _RegistrationsView extends StatelessWidget {
  const _RegistrationsView();

  @override
  Widget build(BuildContext context) {
    return BlocConsumer<RegistrationCubit, RegistrationState>(
      listener: (context, state) {
        if (state is RegistrationLoaded && state.activeFilter == 'pending') {
          sl<PendingCountNotifier>().update(state.items.length);
        }
        if (state is RegistrationEmpty && state.activeFilter == 'pending') {
          sl<PendingCountNotifier>().update(0);
        }
        if (state is RegistrationError) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text(state.message),
              backgroundColor: AppColors.error,
            ),
          );
        }
        if (state is RegistrationActionSuccess) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text(state.message),
              backgroundColor: AppColors.success,
            ),
          );
        }
        if (state is ClientCreateSuccess) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text(state.message),
              backgroundColor: AppColors.success,
            ),
          );
        }
        if (state is ClientCreateError) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text(state.message),
              backgroundColor: AppColors.error,
            ),
          );
        }
      },
      builder: (context, state) {
        final activeFilter = _getActiveFilter(state);
        return Scaffold(
          backgroundColor: AppColors.background,
          appBar: AppBar(
            title: const Text('Pengajuan Customer'),
            actions: [
              IconButton(
                icon: const Icon(Icons.refresh_rounded),
                onPressed: () => context
                    .read<RegistrationCubit>()
                    .loadRegistrations(status: activeFilter),
                tooltip: 'Muat ulang',
              ),
            ],
          ),
          body: Column(
            children: [
              _FilterBar(activeFilter: activeFilter),
              Expanded(child: _buildBody(context, state)),
            ],
          ),
          floatingActionButton: FloatingActionButton.extended(
            onPressed: () => _openCreateClientPage(context),
            backgroundColor: AppColors.primary700,
            foregroundColor: AppColors.white,
            icon: const Icon(Icons.add_business_rounded),
            label: const Text('Client Baru'),
          ),
        );
      },
    );
  }

  void _openCreateClientPage(BuildContext context) {
    final cubit = context.read<RegistrationCubit>();
    Navigator.of(context).push(
      MaterialPageRoute(
        builder: (_) => BlocProvider.value(
          value: cubit,
          child: const CreateClientPage(),
        ),
        fullscreenDialog: true,
      ),
    );
  }

  String? _getActiveFilter(RegistrationState state) {
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
    if (state is RegistrationActionLoading) {
      return _RegistrationList(items: state.items, isLoading: true);
    }
    if (state is ClientCreateLoading) {
      return _RegistrationList(items: state.items, isLoading: true);
    }
    if (state is RegistrationLoaded) {
      return _RegistrationList(items: state.items);
    }
    if (state is RegistrationActionSuccess) {
      return _RegistrationList(items: state.items);
    }
    if (state is ClientCreateSuccess) {
      return _RegistrationList(items: state.items);
    }
    if (state is ClientCreateError) {
      return _RegistrationList(items: state.items);
    }
    if (state is RegistrationEmpty) {
      return _EmptyState(filter: state.activeFilter ?? 'pending');
    }
    if (state is RegistrationError) {
      return Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Icon(Icons.error_outline, size: 48, color: AppColors.error),
            const SizedBox(height: 12),
            Text(state.message,
                style: const TextStyle(color: AppColors.neutral700)),
            const SizedBox(height: 16),
            ElevatedButton(
              onPressed: () =>
                  context.read<RegistrationCubit>().loadRegistrations(),
              child: const Text('Coba Lagi'),
            ),
          ],
        ),
      );
    }
    return const SizedBox.shrink();
  }
}

class _FilterBar extends StatelessWidget {
  final String? activeFilter;
  const _FilterBar({this.activeFilter});

  @override
  Widget build(BuildContext context) {
    final filters = [
      ('pending', 'Menunggu'),
      ('approved', 'Disetujui'),
      ('rejected', 'Ditolak'),
      (null, 'Semua'),
    ];

    return Container(
      color: Colors.white,
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 10),
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
                color: isActive ? AppColors.primary700 : AppColors.neutral700,
                fontWeight:
                    isActive ? FontWeight.w600 : FontWeight.normal,
                fontSize: 13,
              ),
              side: BorderSide(
                color: isActive ? AppColors.primary700 : AppColors.neutral300,
              ),
            ),
          );
        }).toList(),
      ),
    );
  }
}

class _RegistrationList extends StatelessWidget {
  final List<RegistrationEntity> items;
  final bool isLoading;

  const _RegistrationList({required this.items, this.isLoading = false});

  @override
  Widget build(BuildContext context) {
    return Stack(
      children: [
        ListView.separated(
          padding: const EdgeInsets.all(16),
          itemCount: items.length,
          separatorBuilder: (_, __) => const SizedBox(height: 10),
          itemBuilder: (context, i) => _RegistrationCard(item: items[i]),
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

class _RegistrationCard extends StatelessWidget {
  final RegistrationEntity item;
  const _RegistrationCard({required this.item});

  @override
  Widget build(BuildContext context) {
    return Card(
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
                      fontSize: 16,
                      fontWeight: FontWeight.w700,
                      color: AppColors.neutral900,
                    ),
                  ),
                ),
                _StatusBadge(status: item.status),
              ],
            ),
            const SizedBox(height: 10),
            if (item.npwp != null)
              _InfoRow(icon: Icons.receipt_outlined, text: 'NPWP: ${item.npwp}'),
            if (item.email != null)
              _InfoRow(icon: Icons.email_outlined, text: item.email!),
            if (item.telepon != null)
              _InfoRow(icon: Icons.phone_outlined, text: item.telepon!),
            if (item.alamat != null)
              _InfoRow(icon: Icons.location_on_outlined, text: item.alamat!),
            if (item.catatan != null && item.catatan!.isNotEmpty)
              _InfoRow(
                  icon: Icons.notes_outlined,
                  text: 'Catatan: ${item.catatan}'),
            const SizedBox(height: 6),
            Text(
              'Diajukan: ${DateFormat('dd MMM yyyy, HH:mm', 'id_ID').format(item.createdAt.toLocal())}',
              style: const TextStyle(
                fontSize: 12,
                color: AppColors.neutral500,
              ),
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
                      icon: const Icon(Icons.close_rounded, size: 18),
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
                      icon: const Icon(Icons.check_rounded, size: 18),
                      label: const Text('Setujui'),
                      style: ElevatedButton.styleFrom(
                        backgroundColor: AppColors.success,
                        minimumSize: const Size(0, 44),
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
      builder: (dialogCtx) => _ApproveDialog(
        registration: item,
        onApprove: (code, name) {
          context.read<RegistrationCubit>().approve(
                id: item.id,
                companyCode: code,
                companyName: name,
              );
        },
      ),
    );
  }

  void _showRejectDialog(BuildContext context) {
    showDialog(
      context: context,
      builder: (dialogCtx) => _RejectDialog(
        registration: item,
        onReject: (reason) {
          context.read<RegistrationCubit>().reject(
                id: item.id,
                reason: reason,
              );
        },
      ),
    );
  }
}

class _InfoRow extends StatelessWidget {
  final IconData icon;
  final String text;
  const _InfoRow({required this.icon, required this.text});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 4),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Icon(icon, size: 15, color: AppColors.neutral500),
          const SizedBox(width: 6),
          Expanded(
            child: Text(
              text,
              style: const TextStyle(
                fontSize: 13,
                color: AppColors.neutral700,
              ),
            ),
          ),
        ],
      ),
    );
  }
}

class _StatusBadge extends StatelessWidget {
  final String status;
  const _StatusBadge({required this.status});

  @override
  Widget build(BuildContext context) {
    final (label, color, bg) = switch (status) {
      'pending'  => ('Menunggu', AppColors.warning, const Color(0xFFFFF7E0)),
      'approved' => ('Disetujui', AppColors.success, const Color(0xFFE8FFF0)),
      'rejected' => ('Ditolak', AppColors.error, const Color(0xFFFFEEEE)),
      _          => ('Unknown', AppColors.neutral500, AppColors.neutral100),
    };

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
      decoration: BoxDecoration(
        color: bg,
        borderRadius: BorderRadius.circular(20),
      ),
      child: Text(
        label,
        style: TextStyle(
          color: color,
          fontSize: 12,
          fontWeight: FontWeight.w600,
        ),
      ),
    );
  }
}

class _EmptyState extends StatelessWidget {
  final String filter;
  const _EmptyState({required this.filter});

  @override
  Widget build(BuildContext context) {
    final label = switch (filter) {
      'pending'  => 'Tidak ada permintaan yang menunggu',
      'approved' => 'Belum ada registrasi yang disetujui',
      'rejected' => 'Belum ada registrasi yang ditolak',
      _          => 'Belum ada data registrasi',
    };
    return Center(
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          const Icon(Icons.inbox_rounded,
              size: 64, color: AppColors.neutral300),
          const SizedBox(height: 12),
          Text(label,
              style: const TextStyle(
                  color: AppColors.neutral500, fontSize: 15)),
          const SizedBox(height: 16),
          TextButton.icon(
            onPressed: () => context
                .read<RegistrationCubit>()
                .loadRegistrations(status: filter == 'semua' ? null : filter),
            icon: const Icon(Icons.refresh_rounded),
            label: const Text('Muat Ulang'),
          ),
        ],
      ),
    );
  }
}

// ── Dialog Approve ────────────────────────────────────────────────────────────

class _ApproveDialog extends StatefulWidget {
  final RegistrationEntity registration;
  final void Function(String code, String name) onApprove;

  const _ApproveDialog({
    required this.registration,
    required this.onApprove,
  });

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
              'Anda akan menyetujui registrasi dari:\n${widget.registration.nama}',
              style: const TextStyle(fontSize: 14, color: AppColors.neutral700),
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
                helperText: 'Kode unik identifikasi perusahaan',
              ),
              validator: (v) =>
                  v == null || v.trim().isEmpty ? 'Kode wajib diisi' : null,
            ),
            const SizedBox(height: 14),
            TextFormField(
              controller: _nameCtrl,
              decoration: const InputDecoration(
                labelText: 'Nama Resmi Perusahaan *',
                hintText: 'Nama yang akan digunakan di sistem',
              ),
              validator: (v) =>
                  v == null || v.trim().isEmpty ? 'Nama wajib diisi' : null,
            ),
          ],
        ),
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.pop(context),
          child: const Text('Batal'),
        ),
        ElevatedButton(
          onPressed: _submit,
          style: ElevatedButton.styleFrom(
            backgroundColor: AppColors.success,
            minimumSize: const Size(0, 40),
          ),
          child: const Text('Setujui'),
        ),
      ],
    );
  }

  void _submit() {
    if (!_formKey.currentState!.validate()) return;
    Navigator.pop(context);
    widget.onApprove(
      _codeCtrl.text.trim().toUpperCase(),
      _nameCtrl.text.trim(),
    );
  }
}

// ── Dialog Reject ─────────────────────────────────────────────────────────────

class _RejectDialog extends StatefulWidget {
  final RegistrationEntity registration;
  final void Function(String reason) onReject;

  const _RejectDialog({
    required this.registration,
    required this.onReject,
  });

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
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'Anda akan menolak registrasi dari:\n${widget.registration.nama}',
              style: const TextStyle(fontSize: 14, color: AppColors.neutral700),
            ),
            const SizedBox(height: 20),
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
          child: const Text('Batal'),
        ),
        ElevatedButton(
          onPressed: _submit,
          style: ElevatedButton.styleFrom(
            backgroundColor: AppColors.error,
            minimumSize: const Size(0, 40),
          ),
          child: const Text('Tolak'),
        ),
      ],
    );
  }

  void _submit() {
    if (!_formKey.currentState!.validate()) return;
    Navigator.pop(context);
    widget.onReject(_reasonCtrl.text.trim());
  }
}

