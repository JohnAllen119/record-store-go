#include <stdio.h>
#include <stdlib.h>
typedef struct Node
{
    int data;
    struct Node *next;
} Node;

Node *createNode(int value)
{
    Node *newNode = malloc(sizeof(Node));
    if (newNode == NULL)
        return NULL;
    newNode->data = value;
    newNode->next = NULL;
    return newNode;
}
void freeList(Node *head)
{
    Node *current;
    while (head != NULL)
    {
        current = head->next;
        free(head);
        head = current;
    }
}
Node *insertAtHead(Node *head, int value)
{
    Node *newNode = createNode(value);
    if (newNode == NULL)
    {
        return head;
    }
    newNode->next = head;
    return newNode;
}
Node *insertAtTail(Node *head, int value)
{

    Node *newNode = createNode(value);
    if (newNode == NULL)
    {
        return head;
    }
    if (head == NULL)
    {
        return newNode;
    }
    Node *current = head;
    while (current->next)
    {
        current = current->next;
    }
    current->next = newNode;
    return head;
}
Node *findNode(Node *head, int value)
{
    Node *current = head;
    while (current)
    {
        if (current->data == value)
        {
            return current;
        }
        current = current->next;
    }
    return NULL;
}
Node *deleteHead(Node *head)
{
    if (head == NULL)
        return NULL;
    Node *current = head->next;
    free(head);
    head = current;
    return head;
}
Node *deleteByValue(Node *head, int value)
{
    if (head == NULL)
        return NULL;
    if (head->data == value)
    {
        Node *newHead = head->next;
        free(head);
        return newHead;
    }
    Node *current = head;
    while (current->next != NULL && current->next->data != value)
    {
        current = current->next;
    }
    if (current->next != NULL)
    {
        Node *deletedNode = current->next;
        current->next = deletedNode->next;
        free(deletedNode);
    }
    return head;
}

void printList(Node *head)
{
    if (head == NULL)
        return;
    Node *current = head;
    while (current)
    {
        printf("%d ", current->data);
        current = current->next;
    }
}
Node *insertAtIndex(Node *head, int index, int value)
{
    if (index < 0)
    {
        return head;
    }
    if (index == 0)
    {
        head = insertAtHead(head, value);
        return head;
    }
    Node *current = head;
    int cnt = index - 1;
    while (cnt > 0 && current != NULL)
    {
        current = current->next;
        cnt--;
    }
    if (current == NULL)
        return head;
    Node *newNode = createNode(value);
    if (newNode == NULL)
        return head;
    newNode->next = current->next;
    current->next = newNode;
    return head;
}
Node *deleteAtIndex(Node *head, int index)
{
    if (head == NULL || index < 0)
        return head;
    if (index == 0)
    {
        head = deleteHead(head);
        return head;
    }
    Node *current = head;
    int cnt = index - 1;
    while (current != NULL && cnt > 0)
    {
        current = current->next;
        cnt--;
    }
    if (current == NULL ||current->next == NULL )
        return head;
    Node *nxt = current->next;
    current->next = nxt->next;
    free(nxt);
    return head;
}
int main(void)
{
    Node *head = NULL;

    head = insertAtTail(head, 10);
    head = insertAtTail(head, 20);
    head = insertAtTail(head, 30);

    head = insertAtIndex(head, 0, 5);
    head = insertAtIndex(head, 2, 15);
    head = insertAtIndex(head, 5, 40);
    head = insertAtIndex(head, 99, 100);
    head = deleteAtIndex(head, 0);  // 删除5
    head = deleteAtIndex(head, 2);  // 删除20
    head = deleteAtIndex(head, 3);  // 删除40
    head = deleteAtIndex(head, 99); // 不改变链表
    printList(head);
    printf("\n");
    freeList(head);
    return 0;
}