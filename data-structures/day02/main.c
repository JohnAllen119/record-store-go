#include <stdio.h>
#define MAX_SIZE 100
typedef struct
{
    int data[MAX_SIZE];
    int length;
} SeqList;
void initList(SeqList *list)
{
    list->length = 0;
}
int appendElement(SeqList *list, int value)
{
    if (list->length >= MAX_SIZE)
    {
        return 0;
    }
    list->data[list->length] = value;
    list->length++;
    return 1;
}
int insertElement(SeqList *list, int index, int value)
{
    if (list->length >= MAX_SIZE)
    {
        return 0;
    }

    if (index < 0 || index > list->length)
    {
        return 0;
    }

    for (int i = list->length; i > index; i--)
    {
        list->data[i] = list->data[i - 1];
    }
    list->data[index] = value;
    list->length++;
    return 1;
}
void printList(SeqList *list)
{
    for (int i = 0; i < list->length; i++)
    {
        printf("%d ", list->data[i]);
    }
    printf("\n");
}
int deleteElement(SeqList *list, int index, int *deletedValue)
{
    if (0 <= index && index < list->length)
    {
        *deletedValue = list->data[index];
        for (int i = index; i < list->length - 1; i++)
        {
            list->data[i] = list->data[i + 1];
        }
        list->length--;
        return 1;
    }
    return 0;
}
int findElement(SeqList *list, int value)
{
    for (int i = 0; i < list->length; i++)
    {
        if (list->data[i] == value)
        {
            return i;
        }
    }
    return -1;
}

int main(void)
{
    SeqList list;
    initList(&list);
    printf("%d\n", list.length);
     appendElement(&list, 10);
     appendElement(&list, 20);
     appendElement(&list, 30);
     for (int i = 0; i < list.length; i++)
     {
         printf("%d ", list.data[i]);
     }
     printf("\n");
     printf("length=%d ", list.length);
     printf("\n");
     insertElement(&list, 1, 99);
     printList(&list);

     insertElement(&list, list.length, 40);
     printList(&list);

     insertElement(&list, 99, 50);
     printList(&list);
     int deletedValue;

     int deleted = deleteElement(&list, 1, &deletedValue);
     if (deleted)
     {
         printf("%d\n", deletedValue);
         printList(&list);
     }
     
    int index = findElement(&list, 30);

    if (index != -1)
    {
        printf("找到，下标=%d\n", index);
    }
    else
    {
        printf("没有找到\n");
    }
    return 0;
}