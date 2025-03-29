import { Query, Resolver } from '@nestjs/graphql';
import { PersonalShoppingItem } from 'src/graphql/schema';
import { PersonalShoppingService } from './personal-shopping.service';

@Resolver(() => PersonalShoppingItem)
export class PersonalShoppingResolver {
  constructor(private readonly service: PersonalShoppingService) {}

  @Query(() => [PersonalShoppingItem])
  getPersonalShoppingItems(): PersonalShoppingItem[] {
    return this.service.findAll();
  }
}
