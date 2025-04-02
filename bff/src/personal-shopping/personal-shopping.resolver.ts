import { Args, Mutation, Query, Resolver } from '@nestjs/graphql';
import {
  CreatePersonalShoppingItemInput,
  PersonalShoppingItem,
} from 'src/graphql/schema';
import { PersonalShoppingService } from './personal-shopping.service';
import { UseGuards } from '@nestjs/common';
import { AuthGuard } from 'src/auth/auth.guard';

@UseGuards(AuthGuard)
@Resolver(() => PersonalShoppingItem)
export class PersonalShoppingResolver {
  constructor(private readonly service: PersonalShoppingService) {}

  @Query(() => [PersonalShoppingItem])
  getPersonalShoppingItems(): PersonalShoppingItem[] {
    return this.service.findAll();
  }

  @Mutation(() => PersonalShoppingItem)
  createPersonalShoppingItem(
    @Args('input') input: CreatePersonalShoppingItemInput,
  ): PersonalShoppingItem {
    return this.service.create(input);
  }
}
