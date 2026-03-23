import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:intl/intl.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../injection_container.dart' show sl;
import '../../domain/entities/app_release_entity.dart';
import '../../domain/usecases/get_client_installs_usecase.dart';
import '../cubit/app_update_cubit.dart';
import 'create_release_page.dart';
import 'push_update_page.dart';

// Label tampilan per app ID
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

class AppUpdatesPage extends StatefulWidget {
  const AppUpdatesPage({super.key});

  @override
  State<AppUpdatesPage> createState() => _AppUpdatesPageState();
}

class _AppUpdatesPageState extends State<AppUpdatesPage>
    with SingleTickerProviderStateMixin {
  late final AppUpdateCubit _cubit;
  late final TabController _tabController;
  String? _selectedAppId;

  @override
  void initState() {
    super.initState();
    _cubit = sl<AppUpdateCubit>();
    _tabController = TabController(length: 2, vsync: this);
    _cubit.loadReleases();
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return BlocProvider.value(
      value: _cubit,
      child: BlocListener<AppUpdateCubit, AppUpdateState>(
        listener: (context, state) {
          if (state is AppUpdateError) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                content: Text(state.message),
                backgroundColor: AppColors.error,
              ),
            );
          } else if (state is UpdatePushed) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                content: Text(state.message),
                backgroundColor: AppColors.successBase,
              ),
            );
          } else if (state is ReleaseCreated) {
            ScaffoldMessenger.of(context).showSnackBar(
              const SnackBar(
                content: Text('Rilis berhasil dipublikasikan'),
                backgroundColor: AppColors.successBase,
              ),
            );
          }
        },
        child: Scaffold(
          backgroundColor: AppColors.background,
          appBar: AppBar(
            backgroundColor: AppColors.primary800,
            foregroundColor: AppColors.white,
            title: const Text(
              'Manajemen Update App',
              style: TextStyle(fontWeight: FontWeight.bold),
            ),
            actions: [
              IconButton(
                icon: const Icon(Icons.refresh_rounded),
                tooltip: 'Refresh',
                onPressed: () => _cubit.refresh(),
              ),
            ],
            bottom: TabBar(
              controller: _tabController,
              indicatorColor: AppColors.secondary,
              labelColor: AppColors.white,
              unselectedLabelColor: AppColors.primary100,
              tabs: const [
                Tab(text: 'Rilis Versi'),
                Tab(text: 'Status Klien'),
              ],
            ),
          ),
          body: TabBarView(
            controller: _tabController,
            children: [
              _ReleasesTab(
                selectedAppId: _selectedAppId,
                onAppIdChanged: (v) => setState(() {
                  _selectedAppId = v;
                  _cubit.loadReleases(appId: v);
                }),
              ),
              const _ClientStatusTab(),
            ],
          ),
          floatingActionButton: BlocBuilder<AppUpdateCubit, AppUpdateState>(
            builder: (context, state) {
              return FloatingActionButton.extended(
                backgroundColor: AppColors.primary700,
                foregroundColor: AppColors.white,
                icon: const Icon(Icons.add_rounded),
                label: const Text('Rilis Baru'),
                onPressed: () => Navigator.push(
                  context,
                  MaterialPageRoute(
                    fullscreenDialog: true,
                    builder: (_) => BlocProvider.value(
                      value: _cubit,
                      child: const CreateReleasePage(),
                    ),
                  ),
                ),
              );
            },
          ),
        ),
      ),
    );
  }
}

// ─────────────────────────────────────────────────────────────────────────────
// Tab 1: Rilis Versi
// ─────────────────────────────────────────────────────────────────────────────

class _ReleasesTab extends StatelessWidget {
  final String? selectedAppId;
  final ValueChanged<String?> onAppIdChanged;

  const _ReleasesTab({required this.selectedAppId, required this.onAppIdChanged});

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        _AppIdFilterBar(selectedAppId: selectedAppId, onChanged: onAppIdChanged),
        Expanded(
          child: BlocBuilder<AppUpdateCubit, AppUpdateState>(
            builder: (context, state) {
              if (state is AppUpdateLoading || state is ReleaseCreating) {
                return const Center(child: CircularProgressIndicator());
              }

              List<AppReleaseEntity> releases = [];
              if (state is ReleasesLoaded) releases = state.releases;
              if (state is ReleaseCreated) releases = state.releases;
              if (state is PushingUpdate) releases = state.releases;
              if (state is UpdatePushed) releases = state.releases;
              if (state is PushUpdateError) releases = state.releases;

              if (releases.isEmpty) {
                return _EmptyReleases(appId: selectedAppId);
              }

              return RefreshIndicator(
                color: AppColors.primary700,
                onRefresh: () => context.read<AppUpdateCubit>().refresh(),
                child: ListView.builder(
                  padding: const EdgeInsets.fromLTRB(16, 8, 16, 100),
                  itemCount: releases.length,
                  itemBuilder: (_, i) => _ReleaseCard(release: releases[i]),
                ),
              );
            },
          ),
        ),
      ],
    );
  }
}

class _AppIdFilterBar extends StatelessWidget {
  final String? selectedAppId;
  final ValueChanged<String?> onChanged;

  const _AppIdFilterBar({required this.selectedAppId, required this.onChanged});

  @override
  Widget build(BuildContext context) {
    return Container(
      color: AppColors.surface,
      height: 48,
      child: ListView(
        scrollDirection: Axis.horizontal,
        padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
        children: [
          _FilterChip(
            label: 'Semua',
            selected: selectedAppId == null,
            onTap: () => onChanged(null),
          ),
          const SizedBox(width: 8),
          ...List.generate(_appIds.length, (i) {
            final id = _appIds[i];
            return Padding(
              padding: const EdgeInsets.only(right: 8),
              child: _FilterChip(
                label: _appLabels[id] ?? id,
                selected: selectedAppId == id,
                onTap: () => onChanged(id),
              ),
            );
          }),
        ],
      ),
    );
  }
}

class _FilterChip extends StatelessWidget {
  final String label;
  final bool selected;
  final VoidCallback onTap;

  const _FilterChip({required this.label, required this.selected, required this.onTap});

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onTap,
      child: AnimatedContainer(
        duration: const Duration(milliseconds: 200),
        padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 4),
        decoration: BoxDecoration(
          color: selected ? AppColors.primary700 : AppColors.neutral100,
          borderRadius: BorderRadius.circular(20),
        ),
        child: Text(
          label,
          style: TextStyle(
            fontSize: 13,
            fontWeight: selected ? FontWeight.w600 : FontWeight.w400,
            color: selected ? AppColors.white : AppColors.neutral600,
          ),
        ),
      ),
    );
  }
}

class _ReleaseCard extends StatelessWidget {
  final AppReleaseEntity release;

  const _ReleaseCard({required this.release});

  @override
  Widget build(BuildContext context) {
    final fmt = DateFormat('dd MMM yyyy, HH:mm');

    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      elevation: 0,
      color: AppColors.surface,
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                  decoration: BoxDecoration(
                    color: AppColors.primary50,
                    borderRadius: BorderRadius.circular(8),
                  ),
                  child: Text(
                    _appLabels[release.appId] ?? release.appId,
                    style: const TextStyle(
                      fontSize: 11,
                      fontWeight: FontWeight.w600,
                      color: AppColors.primary700,
                    ),
                  ),
                ),
                const Spacer(),
                if (release.isMandatory)
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                    decoration: BoxDecoration(
                      color: AppColors.errorLight,
                      borderRadius: BorderRadius.circular(8),
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
            const SizedBox(height: 10),
            Row(
              children: [
                Text(
                  'v${release.version}',
                  style: const TextStyle(
                    fontSize: 20,
                    fontWeight: FontWeight.bold,
                    color: AppColors.neutral900,
                    fontFamily: 'JetBrainsMono',
                  ),
                ),
                const SizedBox(width: 8),
                Text(
                  '(${release.versionCode})',
                  style: const TextStyle(
                    fontSize: 13,
                    color: AppColors.neutral500,
                    fontFamily: 'JetBrainsMono',
                  ),
                ),
              ],
            ),
            if (release.releaseNotes.isNotEmpty) ...[
              const SizedBox(height: 8),
              Text(
                release.releaseNotes,
                style: const TextStyle(fontSize: 13, color: AppColors.neutral600),
                maxLines: 3,
                overflow: TextOverflow.ellipsis,
              ),
            ],
            const SizedBox(height: 12),
            Row(
              children: [
                Icon(Icons.schedule_rounded, size: 13, color: AppColors.neutral400),
                const SizedBox(width: 4),
                Text(
                  fmt.format(release.createdAt.toLocal()),
                  style: const TextStyle(fontSize: 12, color: AppColors.neutral400),
                ),
                const Spacer(),
                _PushToClientButton(release: release),
              ],
            ),
          ],
        ),
      ),
    );
  }
}

class _PushToClientButton extends StatelessWidget {
  final AppReleaseEntity release;

  const _PushToClientButton({required this.release});

  @override
  Widget build(BuildContext context) {
    return BlocBuilder<AppUpdateCubit, AppUpdateState>(
      builder: (context, state) {
        final loading = state is PushingUpdate;
        return TextButton.icon(
          style: TextButton.styleFrom(
            foregroundColor: AppColors.primary700,
            padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
          ),
          onPressed: loading
              ? null
              : () => Navigator.push(
                    context,
                    MaterialPageRoute(
                      fullscreenDialog: true,
                      builder: (_) => BlocProvider.value(
                        value: context.read<AppUpdateCubit>(),
                        child: PushUpdatePage(release: release),
                      ),
                    ),
                  ),
          icon: const Icon(Icons.send_rounded, size: 16),
          label: const Text('Push ke Klien', style: TextStyle(fontSize: 13)),
        );
      },
    );
  }
}

class _EmptyReleases extends StatelessWidget {
  final String? appId;

  const _EmptyReleases({this.appId});

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(Icons.system_update_outlined, size: 64, color: AppColors.neutral300),
          const SizedBox(height: 16),
          Text(
            appId != null
                ? 'Belum ada rilis untuk ${_appLabels[appId] ?? appId}'
                : 'Belum ada rilis yang dipublikasikan',
            style: const TextStyle(fontSize: 15, color: AppColors.neutral500),
            textAlign: TextAlign.center,
          ),
          const SizedBox(height: 8),
          const Text(
            'Tap tombol + untuk mempublikasikan versi baru',
            style: TextStyle(fontSize: 13, color: AppColors.neutral400),
          ),
        ],
      ),
    );
  }
}

// ─────────────────────────────────────────────────────────────────────────────
// Tab 2: Status Klien
// ─────────────────────────────────────────────────────────────────────────────

class _ClientStatusTab extends StatefulWidget {
  const _ClientStatusTab();

  @override
  State<_ClientStatusTab> createState() => _ClientStatusTabState();
}

class _ClientStatusTabState extends State<_ClientStatusTab> {
  String? _selectedAppId;

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        _AppIdFilterBar(
          selectedAppId: _selectedAppId,
          onChanged: (v) => setState(() => _selectedAppId = v),
        ),
        Expanded(
          child: _selectedAppId == null
              ? _SelectAppPrompt()
              : _AppInstallsList(appId: _selectedAppId!),
        ),
      ],
    );
  }
}

class _SelectAppPrompt extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return const Center(
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(Icons.touch_app_outlined, size: 64, color: AppColors.neutral300),
          SizedBox(height: 16),
          Text(
            'Pilih aplikasi untuk melihat\nstatus instalasi klien',
            style: TextStyle(fontSize: 15, color: AppColors.neutral500),
            textAlign: TextAlign.center,
          ),
        ],
      ),
    );
  }
}

class _AppInstallsList extends StatefulWidget {
  final String appId;

  const _AppInstallsList({required this.appId});

  @override
  State<_AppInstallsList> createState() => _AppInstallsListState();
}

class _AppInstallsListState extends State<_AppInstallsList> {
  bool _loading = false;
  String? _error;
  List<ClientInstallEntity> _installs = [];

  @override
  void initState() {
    super.initState();
    _load();
  }

  @override
  void didUpdateWidget(_AppInstallsList oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.appId != widget.appId) _load();
  }

  Future<void> _load() async {
    setState(() {
      _loading = true;
      _error = null;
    });

    final result = await sl<GetAppInstallsUseCase>()(widget.appId);
    if (!mounted) return;

    result.fold(
      (failure) => setState(() {
        _error = failure.message;
        _loading = false;
      }),
      (items) => setState(() {
        _installs = items;
        _loading = false;
      }),
    );
  }

  @override
  Widget build(BuildContext context) {
    if (_loading) {
      return const Center(child: CircularProgressIndicator());
    }

    if (_error != null) {
      return Center(
        child: Text(_error!, style: const TextStyle(color: AppColors.error)),
      );
    }

    if (_installs.isEmpty) {
      return Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Icon(Icons.devices_outlined, size: 64, color: AppColors.neutral300),
            const SizedBox(height: 16),
            Text(
              'Belum ada data instalasi untuk\n${_appLabels[widget.appId] ?? widget.appId}',
              style: const TextStyle(fontSize: 15, color: AppColors.neutral500),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 8),
            const Text(
              'Data muncul setelah aplikasi pertama kali\ncek update dari server',
              style: TextStyle(fontSize: 13, color: AppColors.neutral400),
              textAlign: TextAlign.center,
            ),
          ],
        ),
      );
    }

    return ListView.builder(
      padding: const EdgeInsets.fromLTRB(16, 8, 16, 80),
      itemCount: _installs.length,
      itemBuilder: (_, i) => _InstallCard(install: _installs[i]),
    );
  }
}

class _InstallCard extends StatelessWidget {
  final ClientInstallEntity install;

  const _InstallCard({required this.install});

  @override
  Widget build(BuildContext context) {
    final hasUpdate = install.needsUpdate;

    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      elevation: 0,
      color: AppColors.surface,
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Expanded(
                  child: Text(
                    install.companyId,
                    style: const TextStyle(
                      fontSize: 12,
                      color: AppColors.neutral400,
                      fontFamily: 'JetBrainsMono',
                    ),
                  ),
                ),
                if (hasUpdate)
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 3),
                    decoration: BoxDecoration(
                      color: install.forceUpdate ? AppColors.errorLight : AppColors.warningLight,
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: Text(
                      install.forceUpdate ? 'WAJIB UPDATE' : 'ADA UPDATE',
                      style: TextStyle(
                        fontSize: 10,
                        fontWeight: FontWeight.w700,
                        color: install.forceUpdate ? AppColors.errorBase : AppColors.warningBase,
                      ),
                    ),
                  ),
              ],
            ),
            const SizedBox(height: 8),
            Row(
              children: [
                _VersionChip(
                  label: 'Terinstal',
                  version: install.installedVersion.isEmpty ? '-' : install.installedVersion,
                  color: AppColors.neutral100,
                  textColor: AppColors.neutral600,
                ),
                if (hasUpdate) ...[
                  const Padding(
                    padding: EdgeInsets.symmetric(horizontal: 8),
                    child: Icon(Icons.arrow_forward_rounded, size: 16, color: AppColors.neutral400),
                  ),
                  _VersionChip(
                    label: 'Target',
                    version: install.targetVersion,
                    color: AppColors.primary50,
                    textColor: AppColors.primary700,
                  ),
                ],
              ],
            ),
            if (install.lastCheckAt != null) ...[
              const SizedBox(height: 8),
              Row(
                children: [
                  const Icon(Icons.history_rounded, size: 13, color: AppColors.neutral400),
                  const SizedBox(width: 4),
                  Text(
                    'Cek terakhir: ${DateFormat('dd MMM, HH:mm').format(install.lastCheckAt!.toLocal())}',
                    style: const TextStyle(fontSize: 12, color: AppColors.neutral400),
                  ),
                ],
              ),
            ],
          ],
        ),
      ),
    );
  }
}

class _VersionChip extends StatelessWidget {
  final String label;
  final String version;
  final Color color;
  final Color textColor;

  const _VersionChip({
    required this.label,
    required this.version,
    required this.color,
    required this.textColor,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 6),
      decoration: BoxDecoration(
        color: color,
        borderRadius: BorderRadius.circular(8),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            label,
            style: TextStyle(fontSize: 10, color: textColor.withOpacity(0.7)),
          ),
          Text(
            version,
            style: TextStyle(
              fontSize: 14,
              fontWeight: FontWeight.bold,
              color: textColor,
              fontFamily: 'JetBrainsMono',
            ),
          ),
        ],
      ),
    );
  }
}
