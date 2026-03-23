import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_dimensions.dart';
import '../../../../injection_container.dart';
import '../cubit/auth_cubit.dart';

class LoginPage extends StatefulWidget {
  const LoginPage({super.key});

  @override
  State<LoginPage> createState() => _LoginPageState();
}

class _LoginPageState extends State<LoginPage> {
  final _formKey = GlobalKey<FormState>();
  final _identifierController = TextEditingController();
  final _passwordController = TextEditingController();
  bool _obscurePassword = true;

  @override
  void dispose() {
    _identifierController.dispose();
    _passwordController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return BlocProvider(
      create: (_) => sl<AuthCubit>(),
      child: Scaffold(
        body: Stack(
          children: [
            // ── Gradient background ──────────────────────────────────────
            Container(
              decoration: const BoxDecoration(
                gradient: LinearGradient(
                  colors: [AppColors.primary900, AppColors.primary700],
                  begin: Alignment.topLeft,
                  end: Alignment.bottomRight,
                ),
              ),
            ),
            // ── Decorative circles ───────────────────────────────────────
            Positioned(
              top: -60,
              right: -60,
              child: Container(
                width: 220,
                height: 220,
                decoration: BoxDecoration(
                  shape: BoxShape.circle,
                  color: AppColors.textOnPrimary.withValues(alpha: 0.05),
                ),
              ),
            ),
            Positioned(
              top: 80,
              right: 40,
              child: Container(
                width: 100,
                height: 100,
                decoration: BoxDecoration(
                  shape: BoxShape.circle,
                  color: AppColors.secondary400.withValues(alpha: 0.15),
                ),
              ),
            ),
            Positioned(
              bottom: 200,
              left: -40,
              child: Container(
                width: 160,
                height: 160,
                decoration: BoxDecoration(
                  shape: BoxShape.circle,
                  color: AppColors.primary500.withValues(alpha: 0.3),
                ),
              ),
            ),
            // ── Content ──────────────────────────────────────────────────
            SafeArea(
              child: Center(
                child: ConstrainedBox(
                  constraints: const BoxConstraints(maxWidth: 480),
                  child: SingleChildScrollView(
                    padding: const EdgeInsets.symmetric(
                      horizontal: AppDimensions.spacing24,
                      vertical: AppDimensions.spacing32,
                    ),
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        // Logo area
                        Container(
                          width: 72,
                          height: 72,
                          decoration: BoxDecoration(
                            color: AppColors.textOnPrimary.withValues(alpha: 0.15),
                            borderRadius:
                                BorderRadius.circular(AppDimensions.radiusXl),
                            border: Border.all(
                              color:
                                  AppColors.textOnPrimary.withValues(alpha: 0.2),
                            ),
                          ),
                          child: const Icon(
                            Icons.developer_mode_rounded,
                            size: 40,
                            color: AppColors.secondary400,
                          ),
                        ),
                        const SizedBox(height: AppDimensions.spacing16),
                        const Text(
                          'FlashERP Developer',
                          style: TextStyle(
                            fontSize: 28,
                            fontWeight: FontWeight.w800,
                            color: AppColors.textOnPrimary,
                            letterSpacing: 0.5,
                          ),
                        ),
                        const SizedBox(height: AppDimensions.spacing4),
                        Text(
                          'Portal Developer Sales',
                          style: TextStyle(
                            fontSize: 14,
                            color: AppColors.textOnPrimary.withValues(alpha: 0.7),
                            fontWeight: FontWeight.w400,
                          ),
                        ),
                        const SizedBox(height: AppDimensions.spacing32),
                        // Form card
                        Container(
                          width: double.infinity,
                          decoration: const BoxDecoration(
                            color: AppColors.surface,
                            borderRadius: BorderRadius.all(
                              Radius.circular(AppDimensions.radius2xl),
                            ),
                          ),
                          child: Padding(
                            padding: const EdgeInsets.all(AppDimensions.spacing32),
                            child: BlocConsumer<AuthCubit, AuthState>(
                              listener: (context, state) {
                                if (state is AuthAuthenticated) {
                                  // GoRouter redirect akan handle navigasi ke /home
                                }
                                if (state is AuthError) {
                                  ScaffoldMessenger.of(context).showSnackBar(
                                    SnackBar(
                                      content: Text(state.message),
                                      backgroundColor: AppColors.errorBase,
                                      behavior: SnackBarBehavior.floating,
                                    ),
                                  );
                                }
                              },
                              builder: (context, state) => Form(
                                key: _formKey,
                                child: Column(
                                  crossAxisAlignment:
                                      CrossAxisAlignment.stretch,
                                  children: [
                                    const Text(
                                      'Masuk',
                                      style: TextStyle(
                                        fontSize: 22,
                                        fontWeight: FontWeight.w800,
                                        color: AppColors.textPrimary,
                                      ),
                                    ),
                                    const SizedBox(height: AppDimensions.spacing4),
                                    const Text(
                                      'Selamat datang kembali',
                                      style: TextStyle(
                                        fontSize: 14,
                                        color: AppColors.textSecondary,
                                      ),
                                    ),
                                    const SizedBox(
                                        height: AppDimensions.spacing24),
                                    // Email field
                                    TextFormField(
                                      controller: _identifierController,
                                      decoration: InputDecoration(
                                        labelText: 'Email',
                                        hintText: 'nama@flashlab.id',
                                        prefixIcon: const Icon(
                                            Icons.email_outlined),
                                        border: OutlineInputBorder(
                                          borderRadius: BorderRadius.circular(
                                              AppDimensions.radiusMd),
                                        ),
                                      ),
                                      keyboardType:
                                          TextInputType.emailAddress,
                                      style: const TextStyle(
                                          color: AppColors.textPrimary),
                                      validator: (v) {
                                        if (v == null || v.trim().isEmpty) {
                                          return 'Email wajib diisi';
                                        }
                                        return null;
                                      },
                                      textInputAction: TextInputAction.next,
                                    ),
                                    const SizedBox(
                                        height: AppDimensions.spacing16),
                                    // Password field
                                    TextFormField(
                                      controller: _passwordController,
                                      decoration: InputDecoration(
                                        labelText: 'Password',
                                        prefixIcon: const Icon(
                                            Icons.lock_outline),
                                        suffixIcon: IconButton(
                                          icon: Icon(
                                            _obscurePassword
                                                ? Icons.visibility_off_outlined
                                                : Icons.visibility_outlined,
                                          ),
                                          onPressed: () => setState(() =>
                                              _obscurePassword =
                                                  !_obscurePassword),
                                        ),
                                        border: OutlineInputBorder(
                                          borderRadius: BorderRadius.circular(
                                              AppDimensions.radiusMd),
                                        ),
                                      ),
                                      style: const TextStyle(
                                          color: AppColors.textPrimary),
                                      obscureText: _obscurePassword,
                                      validator: (v) =>
                                          v == null || v.isEmpty
                                              ? 'Password wajib diisi'
                                              : null,
                                      textInputAction: TextInputAction.done,
                                      onFieldSubmitted: (_) =>
                                          _submit(context, state),
                                    ),
                                    const SizedBox(
                                        height: AppDimensions.spacing24),
                                    // Login button
                                    FilledButton(
                                      onPressed: state is AuthLoading
                                          ? null
                                          : () => _submit(context, state),
                                      style: FilledButton.styleFrom(
                                        minimumSize:
                                            const Size.fromHeight(52),
                                        backgroundColor: AppColors.primary600,
                                        shape: RoundedRectangleBorder(
                                          borderRadius:
                                              BorderRadius.circular(
                                                  AppDimensions.radiusMd),
                                        ),
                                      ),
                                      child: state is AuthLoading
                                          ? const SizedBox(
                                              height: 20,
                                              width: 20,
                                              child:
                                                  CircularProgressIndicator(
                                                strokeWidth: 2,
                                                color: Colors.white,
                                              ),
                                            )
                                          : const Text(
                                              'Masuk',
                                              style: TextStyle(
                                                fontSize: 16,
                                                fontWeight: FontWeight.w700,
                                              ),
                                            ),
                                    ),
                                  ],
                                ),
                              ),
                            ),
                          ),
                        ),
                        const SizedBox(height: AppDimensions.spacing20),
                        // Footer
                        RichText(
                          text: TextSpan(
                            style: TextStyle(
                              fontSize: 12,
                              color: AppColors.textOnPrimary
                                  .withValues(alpha: 0.6),
                            ),
                            children: [
                              const TextSpan(text: 'Powered by '),
                              TextSpan(
                                text: 'Thunderlab',
                                style: TextStyle(
                                  color: AppColors.secondary400,
                                  fontWeight: FontWeight.w600,
                                ),
                              ),
                            ],
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  void _submit(BuildContext context, AuthState state) {
    if (state is AuthLoading) return;
    if (_formKey.currentState?.validate() != true) return;
    context.read<AuthCubit>().login(
          identifier: _identifierController.text.trim(),
          password: _passwordController.text,
        );
  }
}
