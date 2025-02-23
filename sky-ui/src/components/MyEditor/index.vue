<template>
  <div ref="editor" class="editor"></div>
</template>

<script setup>
import { onMounted, ref, watch, defineProps, defineEmits } from 'vue';
import Quill from 'quill';
import 'quill/dist/quill.snow.css';

const props = defineProps({
  modelValue: String,
  minHeight: Number
});
const emit = defineEmits(['update:modelValue']);

const editor = ref(null);
const quillInstance = ref(null);

onMounted(() => {
  quillInstance.value = new Quill(editor.value, {
    theme: 'snow',
    placeholder: '请输入内容...',
    modules: {
      toolbar: [
        [{ header: [1, 2, false] }],
        ['bold', 'italic', 'underline'],
        ['image', 'code-block'],
        [{ list: 'ordered' }, { list: 'bullet' }],
        [{ script: 'sub' }, { script: 'super' }],
        [{ indent: '-1' }, { indent: '+1' }],
        ['clean']
      ]
    }
  });

  quillInstance.value.root.innerHTML = props.modelValue;

  quillInstance.value.on('text-change', () => {
    emit('update:modelValue', quillInstance.value.root.innerHTML);
  });
});

watch(() => props.modelValue, (newVal) => {
  if (quillInstance.value.root.innerHTML !== newVal) {
    quillInstance.value.root.innerHTML = newVal;
  }
});

</script>

<style scoped>
.editor {
  min-height: v-bind('minHeight + "px"');
  /* background-color: #fff; 设置编辑器背景为白色 */
  color: #333; /* 设置编辑器文字颜色为黑色 */
  border: 1px solid #ccc; /* 设置编辑器边框颜色为浅灰色 */
}

::v-deep .ql-container {
  /* background-color: #fff; 设置编辑器内容区域背景为白色 */
  color: #333; /* 设置编辑器内容区域文字颜色为黑色 */
}

::v-deep .ql-toolbar {
  /* background-color: #fff; 设置工具栏背景为白色 */
  color: #333; /* 设置工具栏文字颜色为黑色 */
}

::v-deep .ql-picker-label {
  color: #333; /* 设置下拉菜单标签文字颜色为黑色 */
}

::v-deep .ql-picker-item {
  color: #333; /* 设置下拉菜单项文字颜色为黑色 */
}

::v-deep .ql-editor {
  height: v-bind('minHeight + "px"');
  /* min-height: v-bind('minHeight + "px"'); 设置编辑器内容区域最小高度 */
  /* background-color: #fff; 设置编辑器内容区域背景为白色 */
  color: #333; /* 设置编辑器内容区域文字颜色为黑色 */
}
</style>