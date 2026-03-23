import 'package:flutter_bloc/flutter_bloc.dart';

import '../../../../core/auth/auth_notifier.dart';
import '../../domain/entities/user_entity.dart';
import '../../domain/usecases/login_usecase.dart';

part 'auth_state.dart';

class AuthCubit extends Cubit<AuthState> {
  final LoginUseCase _loginUseCase;
  final AuthNotifier _authNotifier;

  AuthCubit({
    required LoginUseCase loginUseCase,
    required AuthNotifier authNotifier,
  })  : _loginUseCase = loginUseCase,
        _authNotifier = authNotifier,
        super(AuthInitial());

  Future<void> login({
    required String identifier,
    required String password,
  }) async {
    emit(AuthLoading());
    final result = await _loginUseCase(
      identifier: identifier,
      password: password,
    );
    if (result.isLeft()) {
      final failure = result.fold((f) => f, (_) => null)!;
      emit(AuthError(failure.message));
      return;
    }
    final user = result.getOrElse(() => throw Exception());
    await _authNotifier.onLogin(user.role, user.id);
    emit(AuthAuthenticated(user));
  }

  Future<void> logout() async {
    await _authNotifier.onLogout();
    emit(AuthUnauthenticated());
  }
}
