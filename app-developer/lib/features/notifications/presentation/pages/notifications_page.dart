import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:intl/intl.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_dimensions.dart';
import '../../../../core/notifiers/pending_count_notifier.dart';
import '../../../../injection_container.dart';
import '../../../registrations/domain/entities/registration_entity.dart';
import '../../../registrations/presentation/cubit/registration_cubit.dart';

class NotificationsPage extends StatelessWidget {
  const NotificationsPage({super.key});

  @override
  Widget build(BuildContext context) {
    return BlocProvider(
      create: (_) =>
          sl<RegistrationCubit>()..loadRegistrations(status: 'pending'),
      child: const _NotificationsView(),
    );
  }
}

class _NotificationsView extends StatelessWidget {
  const _NotificationsView();

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
        if (state is RegistrationActionSuccess) {
          ScaffoldMessenger.of(context).showSnackBar(SnackBar(
            content: Text(state.message),
            backgroundColor: AppColors.success,
          ));
          context
              .read<RegistrationCubit>()
              .loadRegistrations(status: 'pending');
        }
        if (state is RegistrationError) {
          ScaffoldMessenger.of(context).showSnackBar(SnackBar(
            content: Text(state.message),
            backgroundColor: AppColors.error,
          ));
        }
      },
      builder: (context, state) => Scaffold(
        backgroundColor: AppColors.background,
        appBar: AppBar(
          backgroundColor: AppColors.surface,
          title: const Text('Notifikasi'),
          actions: [
            IconButton(
              icon: const Icon(Icons.refresh_rounded),
              tooltip: 'Muat ulang',
              onPressed: () => context
                  .read<RegistrationCubit>()
                  .loadRegistrations(status: 'pending'),
            ),
          ],
        ),
        body: _buildBody(context, state),
      ),
    );
  }

  Widget _buildBody(BuildContext context, RegistrationState state) {
    if (state is RegistrationLoading || state is RegistrationInitial) {
      return const Center(child: CircularProgressIndicator());
    }
    if (state is RegistrationActionLoading) {
      return _TaskList(items: state.items, isLoading: true);
    }
    if (state is RegistrationLoaded || state is RegistrationActionSuccess) {
      final items = state is RegistrationLoaded
          ? state.items
          : (state as RegistrationActionSuccess).items;
      return _TaskList(items: items);
    }
    if (state is RegistrationEmpty) {
      return const _EmptyNotifications();
    }
    if (state is RegistrationError) {
      return Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Icon(Icons.error_outline,
                size: 48, color: AppColors.errorBase),
            const SizedBox(height: AppDimensions.spacing12),
            Text(state.message,
                textAlign: TextAlign.center,
                style: const TextStyle(color: AppColors.neutral700)),
            const SizedBox(height: AppDimensions.spacing16),
            FilledButton(
              onPressed: () => context
                  .read<RegistrationCubit>()
                  .loadRegistrations(status: 'pending'),
              child: const Text('Coba Lagi'),
            ),
          ],
        ),
      );
    }
    return const SizedBox.shrink();
  }
}

// ── Task List ──────────────────────────────────────────────────────────────────

class _TaskList extends StatelessWidget {
  final List<RegistrationEntity> items;
  final bool isLoading;

  const _TaskList({required this.items, this.isLoading = false});

  @override
  Widget build(BuildContext context) {
    return Stack(
      children: [
        ListView(
          padding: const EdgeInsets.symmetric(
            horizontal: AppDimensions.screenPaddingH,
            vertical: AppDimensions.screenPaddingV,
          ),
          children: [
            _SectionHeader(count: items.length),
            const SizedBox(height: AppDimensions.spacing12),
            ...items.map((item) => _TaskCard(item: item)),
            const SizedBox(height: AppDimensions.spacing80),
          ],
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

class _SectionHeader extends StatelessWidget {
  final int count;
  const _SectionHeader({required this.count});

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        Container(
          width: 4,
          height: 20,
          decoration: BoxDecoration(
            color: AppColors.warningBase,
            borderRadius: BorderRadius.circular(AppDimensions.radiusFull),
          ),
        ),
        const SizedBox(width: AppDimensions.spacing8),
        Text(
          'Pengajuan Menunggu Persetujuan',
          style: Theme.of(context)
              .textTheme
              .titleSmall
              ?.copyWith(color: AppColors.neutral700),
        ),
        const SizedBox(width: AppDimensions.spacing8),
        Container(
          padding: const EdgeInsets.symmetric(
            horizontal: AppDimensions.spacing8,
            vertical: AppDimensions.spacing2,
          ),
          decoration: BoxDecoration(
            color: AppColors.warningLight,
            borderRadius: BorderRadius.circular(AppDimensions.radiusFull),
          ),
          child: Text(
            '$count',
            style: const TextStyle(
              color: AppColors.warningBase,
              fontSize: 12,
              fontWeight: FontWeight.w700,
            ),
          ),
        ),
      ],
    );
  }
}

// ── Task Card ──────────────────────────────────────────────────────────────────

class _TaskCard extends StatelessWidget {
  final RegistrationEntity item;
  const _TaskCard({required this.item});

  @override
  Widget build(BuildContext context) {
    return Container(
      margin: const EdgeInsets.only(bottom: AppDimensions.spacing12),
      decoration: BoxDecoration(
        color: AppColors.surface,
        borderRadius: BorderRadius.circular(AppDimensions.radiusXl),
        border: Border.all(color: AppColors.neutral200),
        boxShadow: const [
          BoxShadow(
            color: Color(0x08000000),
            blurRadius: 8,
            offset: Offset(0, 2),
          ),
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // ── Header ──
          Padding(
            padding: const EdgeInsets.fromLTRB(
              AppDimensions.cardPadding,
              AppDimensions.cardPadding,
              AppDimensions.cardPadding,
              AppDimensions.spacing8,
            ),
            child: Row(
              children: [
                Container(
                  width: 44,
                  height: 44,
                  decoration: BoxDecoration(
                    color: AppColors.primary50,
                    borderRadius:
                        BorderRadius.circular(AppDimensions.radiusMd),
                  ),
                  child: const Icon(Icons.business_rounded,
                      color: AppColors.primary700,
                      size: AppDimensions.iconMd),
                ),
                const SizedBox(width: AppDimensions.spacing12),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        item.nama,
                        style: Theme.of(context)
                            .textTheme
                            .titleSmall
                            ?.copyWith(color: AppColors.neutral900),
                        maxLines: 1,
                        overflow: TextOverflow.ellipsis,
                      ),
                      const SizedBox(height: AppDimensions.spacing2),
                      Text(
                        DateFormat('dd MMM yyyy, HH:mm', 'id_ID')
                            .format(item.createdAt.toLocal()),
                        style: const TextStyle(
                          fontSize: 12,
                          color: AppColors.neutral500,
                        ),
                      ),
                    ],
                  ),
                ),
                Container(
                  padding: const EdgeInsets.symmetric(
                    horizontal: AppDimensions.spacing8,
                    vertical: AppDimensions.spacing4,
                  ),
                  decoration: BoxDecoration(
                    color: AppColors.warningLight,
                    borderRadius:
                        BorderRadius.circular(AppDimensions.radiusFull),
                  ),
                  child: const Text(
                    'Menunggu',
                    style: TextStyle(
                      color: AppColors.warningBase,
                      fontSize: 11,
                      fontWeight: FontWeight.w700,
                    ),
                  ),
                ),
              ],
            ),
          ),
          // ── Detail ──
          if (item.email != null ||
              item.telepon != null ||
              item.npwp != null ||
              item.alamat != null)
            Padding(
              padding: const EdgeInsets.symmetric(
                  horizontal: AppDimensions.cardPadding),
              child: Column(
                children: [
                  if (item.email != null)
                    _InfoRow(Icons.email_outlined, item.email!),
                  if (item.telepon != null)
                    _InfoRow(Icons.phone_outlined, item.telepon!),
                  if (item.npwp != null)
                    _InfoRow(Icons.receipt_outlined, 'NPWP: ${item.npwp}'),
                  if (item.alamat != null)
                    _InfoRow(Icons.location_on_outlined, item.alamat!),
                  const SizedBox(height: AppDimensions.spacing8),
                ],
              ),
            ),
          // ── Actions ──
          Padding(
            padding: const EdgeInsets.all(AppDimensions.cardPadding),
            child: Row(
              children: [
                Expanded(
                  child: OutlinedButton.icon(
                    onPressed: () => _showRejectDialog(context),
                    icon: const Icon(Icons.close_rounded, size: 18),
                    label: const Text('Tolak'),
                    style: OutlinedButton.styleFrom(
                      foregroundColor: AppColors.errorBase,
                      side: const BorderSide(color: AppColors.errorBase),
                    ),
                  ),
                ),
                const SizedBox(width: AppDimensions.spacing12),
                Expanded(
                  child: FilledButton.icon(
                    onPressed: () => _showApproveDialog(context),
                    icon: const Icon(Icons.check_rounded, size: 18),
                    label: const Text('Setujui'),
                    style: FilledButton.styleFrom(
                      backgroundColor: AppColors.successBase,
                    ),
                  ),
                ),
              ],
            ),
          ),
        ],
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
  const _InfoRow(this.icon, this.text);

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(bottom: AppDimensions.spacing4),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Icon(icon, size: 14, color: AppColors.neutral400),
          const SizedBox(width: AppDimensions.spacing8),
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

// ── Empty State ────────────────────────────────────────────────────────────────

class _EmptyNotifications extends StatelessWidget {
  const _EmptyNotifications();

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Container(
            width: 80,
            height: 80,
            decoration: const BoxDecoration(
              color: AppColors.successLight,
              shape: BoxShape.circle,
            ),
            child: const Icon(
              Icons.check_circle_outline_rounded,
              size: 40,
              color: AppColors.successBase,
            ),
          ),
          const SizedBox(height: AppDimensions.spacing16),
          Text(
            'Semua Beres!',
            style: Theme.of(context)
                .textTheme
                .titleMedium
                ?.copyWith(color: AppColors.neutral900),
          ),
          const SizedBox(height: AppDimensions.spacing8),
          const Text(
            'Tidak ada pengajuan yang perlu diproses.',
            style: TextStyle(color: AppColors.neutral500, fontSize: 14),
          ),
        ],
      ),
    );
  }
}

// ── Dialogs ────────────────────────────────────────────────────────────────────

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
      title: const Text('Setujui Pengajuan'),
      content: Form(
        key: _formKey,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'Anda akan menyetujui pengajuan dari:\n${widget.registration.nama}',
              style:
                  const TextStyle(fontSize: 14, color: AppColors.neutral700),
            ),
            const SizedBox(height: AppDimensions.spacing4),
            const Text(
              'Tindakan ini akan membuat perusahaan baru di sistem.',
              style: TextStyle(fontSize: 12, color: AppColors.neutral500),
            ),
            const SizedBox(height: AppDimensions.spacing20),
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
            const SizedBox(height: AppDimensions.spacing16),
            TextFormField(
              controller: _nameCtrl,
              decoration: const InputDecoration(
                labelText: 'Nama Resmi Perusahaan *',
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
        FilledButton(
          onPressed: _submit,
          style:
              FilledButton.styleFrom(backgroundColor: AppColors.successBase),
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
      title: const Text('Tolak Pengajuan'),
      content: Form(
        key: _formKey,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'Anda akan menolak pengajuan dari:\n${widget.registration.nama}',
              style:
                  const TextStyle(fontSize: 14, color: AppColors.neutral700),
            ),
            const SizedBox(height: AppDimensions.spacing20),
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
        FilledButton(
          onPressed: _submit,
          style: FilledButton.styleFrom(backgroundColor: AppColors.errorBase),
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
