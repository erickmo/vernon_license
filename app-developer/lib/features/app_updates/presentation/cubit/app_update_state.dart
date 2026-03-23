part of 'app_update_cubit.dart';

sealed class AppUpdateState {}

class AppUpdateInitial extends AppUpdateState {}

class AppUpdateLoading extends AppUpdateState {}

class ReleasesLoaded extends AppUpdateState {
  final List<AppReleaseEntity> releases;
  final String? appIdFilter;
  ReleasesLoaded(this.releases, {this.appIdFilter});
}

class ReleasesEmpty extends AppUpdateState {
  final String? appIdFilter;
  ReleasesEmpty({this.appIdFilter});
}

class AppUpdateError extends AppUpdateState {
  final String message;
  AppUpdateError(this.message);
}

class ReleaseCreating extends AppUpdateState {}

class ReleaseCreated extends AppUpdateState {
  final List<AppReleaseEntity> releases;
  ReleaseCreated(this.releases);
}

class PushingUpdate extends AppUpdateState {
  final List<AppReleaseEntity> releases;
  PushingUpdate(this.releases);
}

class UpdatePushed extends AppUpdateState {
  final String message;
  final List<AppReleaseEntity> releases;
  UpdatePushed(this.message, this.releases);
}

class PushUpdateError extends AppUpdateState {
  final String message;
  final List<AppReleaseEntity> releases;
  PushUpdateError(this.message, this.releases);
}
