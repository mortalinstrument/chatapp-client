<template>
  <ion-page>
    <ion-header :translucent="true">
      <ion-toolbar>
        <ion-title>SchardtChat</ion-title>
        <ion-button :collapse="true" slot="end">
          <ion-icon slot="icon-only" :icon="cog"></ion-icon>
        </ion-button>
      </ion-toolbar>
    </ion-header>
    <ion-content :fullscreen="true">
      <ion-header collapse="condense">
        <ion-toolbar>
        <ion-title>SchardtChat</ion-title>
        <ion-button :collapse="true" slot="end">
          <ion-icon slot="icon-only" :icon="cog"></ion-icon>
        </ion-button>
      </ion-toolbar>
    </ion-header>

    <ListCardVue :partners="getPartners()"></ListCardVue>
      
  </ion-content>
  </ion-page>
</template>

<script lang="ts">
import { IonContent, IonHeader, IonPage, IonTitle, IonToolbar, IonButton, IonIcon } from '@ionic/vue';
import { star, cog } from 'ionicons/icons';
import { defineComponent } from 'vue';
import { mapGetters, mapActions } from 'vuex';
import ListCardVue from '@/components/ListCard.vue';

export default defineComponent({
  name: 'HomePage',
  components: {
    IonContent,
    IonHeader,
    IonPage,
    IonTitle,
    IonToolbar,
    IonIcon,
    IonButton,
    ListCardVue
  },
  beforeMount(){
    this.initWebsocket()
  },
  setup() {
    return { star, cog }
  },
  methods: {
  ...mapGetters([
    'getPartners'
  ]),
  ...mapActions({
    initWebsocket: 'establishWebsocketConnection'
  }),
  }
});
</script>

<style scoped>
#container {
  text-align: center;
  
  position: absolute;
  left: 0;
  right: 0;
  top: 50%;
  transform: translateY(-50%);
}

#container strong {
  font-size: 20px;
  line-height: 26px;
}

#container p {
  font-size: 16px;
  line-height: 22px;
  
  color: #8c8c8c;
  
  margin: 0;
}

#container a {
  text-decoration: none;
}
</style>
