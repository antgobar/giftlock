<template>
  <section class="section">
    <div class="container">
      <div class="columns is-centered">
        <div class="column is-8">
          <div class="box">
            <h1 class="title is-2 has-text-centered">Welcome to Gift Lock!</h1>

            <div class="buttons is-centered">
              <form action="">
                <div class="field has-addons">
                  <div class="control">
                    <input class="input" type="text" placeholder="Create a new gift group" />
                  </div>
                  <div class="control">
                    <button class="button is-primary">
                      Create
                    </button>
                  </div>
                </div>
              </form>
              
              <button class="button is-info">
                My Gifts
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
import { useRouter } from 'vue-router';

const router = useRouter()

</script>

