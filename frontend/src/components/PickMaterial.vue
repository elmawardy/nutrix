<template>
    <div>
        <DataTable  filterDisplay="row" :loading="loading" v-model:filters="filters" :globalFilterFields="['name']" :value="materials" stripedRows tableStyle="min-width: 50rem" class="w-full pr-5">
            <template #header>
                <div class="flex justify-content-start">
                    <IconField iconPosition="left">
                        <InputIcon>
                            <i class="pi pi-search" />
                        </InputIcon>
                        <InputText v-model="filters['name'].value" placeholder="Search by name" />
                    </IconField>
                </div>
            </template>
            <Column field="name" header="Name"></Column>
            <Column field="quantity" header="Quantity"></Column>
            <Column field="unit" header="Unit"></Column>
            <Column header="Actions">
                <template #body="slotProps">
                    <ButtonGroup>
                        <Button icon="pi pi-plus" label="Add" severity="secondary" aria-label="Ddd" @click="returnMaterial(slotProps.data)" />
                    </ButtonGroup>
                </template>
            </Column>
        </DataTable>
    </div>
</template>

<script setup lang="ts">
import {ref,defineEmits} from 'vue'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext';
import InputIcon from 'primevue/inputicon'
import IconField from 'primevue/iconfield'
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import axios from 'axios'
import { FilterMatchMode } from 'primevue/api';
import { Material, MaterialEntry } from '@/classes/OrderItem';


const materials = ref([])
const loading = ref(false)

const filters = ref({
    name: { value: null, matchMode: FilterMatchMode.CONTAINS },
});


const emit = defineEmits(['returnMaterial'])


const returnMaterial = (material: Material) => {
    emit('returnMaterial', material)
}


const GetMaterials = () => {
    loading.value = true
    axios.get("http://localhost:8000/api/components")
    .then((response) => {

        
        response.data.forEach((component:Material,materialIndex: number) => {

            component.entries?.forEach((entry: MaterialEntry,entryIndex: number) => {
                if (entry.quantity < 0){
                    response.data[materialIndex].entries.splice(entryIndex, 1)
                    return
                }
                    
                component.quantity += entry.quantity
            });
        });


        materials.value = response.data
        loading.value = false
    })
}

GetMaterials()


</script>