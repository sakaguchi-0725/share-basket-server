import { Module } from '@nestjs/common';
import { PersonalShoppingResolver } from './personal-shopping.resolver';
import { PersonalShoppingService } from './personal-shopping.service';

@Module({
  providers: [PersonalShoppingResolver, PersonalShoppingService],
})
export class PersonalShoppingModule {}
