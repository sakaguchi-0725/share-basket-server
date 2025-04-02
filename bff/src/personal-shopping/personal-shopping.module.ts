import { Module } from '@nestjs/common';
import { PersonalShoppingResolver } from './personal-shopping.resolver';
import { PersonalShoppingService } from './personal-shopping.service';
import { AuthModule } from 'src/auth/auth.module';

@Module({
  imports: [AuthModule],
  providers: [PersonalShoppingResolver, PersonalShoppingService],
})
export class PersonalShoppingModule {}
