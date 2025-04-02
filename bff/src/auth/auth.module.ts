import { Module } from '@nestjs/common';
import { AuthResolver } from './auth.resolver';
import { AuthService } from './auth.service';
import { AuthGuard } from './auth.guard';

@Module({
  providers: [AuthResolver, AuthService, AuthGuard],
  exports: [AuthService, AuthGuard],
})
export class AuthModule {}
