import 'package:flutter_bloc/flutter_bloc.dart';

import '../../domain/entities/app_release_entity.dart';
import '../../domain/usecases/create_release_usecase.dart';
import '../../domain/usecases/list_releases_usecase.dart';
import '../../domain/usecases/push_update_usecase.dart';

part 'app_update_state.dart';

class AppUpdateCubit extends Cubit<AppUpdateState> {
  final ListReleasesUseCase _listReleases;
  final CreateReleaseUseCase _createRelease;
  final PushUpdateUseCase _pushUpdate;

  List<AppReleaseEntity> _currentReleases = [];
  String? _currentAppIdFilter;

  AppUpdateCubit({
    required ListReleasesUseCase listReleases,
    required CreateReleaseUseCase createRelease,
    required PushUpdateUseCase pushUpdate,
  })  : _listReleases = listReleases,
        _createRelease = createRelease,
        _pushUpdate = pushUpdate,
        super(AppUpdateInitial());

  Future<void> loadReleases({String? appId}) async {
    _currentAppIdFilter = appId;
    emit(AppUpdateLoading());

    final result = await _listReleases(appId: appId);
    result.fold(
      (failure) => emit(AppUpdateError(failure.message)),
      (releases) {
        _currentReleases = releases;
        if (releases.isEmpty) {
          emit(ReleasesEmpty(appIdFilter: appId));
        } else {
          emit(ReleasesLoaded(releases, appIdFilter: appId));
        }
      },
    );
  }

  Future<void> createRelease({
    required String appId,
    required String version,
    required int versionCode,
    required String downloadUrl,
    String? releaseNotes,
    bool isMandatory = false,
  }) async {
    emit(ReleaseCreating());

    final result = await _createRelease(
      appId: appId,
      version: version,
      versionCode: versionCode,
      downloadUrl: downloadUrl,
      releaseNotes: releaseNotes,
      isMandatory: isMandatory,
    );

    await result.fold(
      (failure) async => emit(AppUpdateError(failure.message)),
      (_) async {
        // Reload setelah berhasil
        final listResult = await _listReleases(appId: _currentAppIdFilter);
        listResult.fold(
          (failure) => emit(AppUpdateError(failure.message)),
          (releases) {
            _currentReleases = releases;
            emit(ReleaseCreated(releases));
          },
        );
      },
    );
  }

  Future<void> pushUpdate({
    required String companyId,
    required String appId,
    required int versionCode,
    bool forceUpdate = false,
  }) async {
    emit(PushingUpdate(_currentReleases));

    final result = await _pushUpdate(
      companyId: companyId,
      appId: appId,
      versionCode: versionCode,
      forceUpdate: forceUpdate,
    );

    result.fold(
      (failure) => emit(PushUpdateError(failure.message, _currentReleases)),
      (_) => emit(UpdatePushed('Update berhasil didorong ke klien', _currentReleases)),
    );
  }

  Future<void> refresh() => loadReleases(appId: _currentAppIdFilter);
}
